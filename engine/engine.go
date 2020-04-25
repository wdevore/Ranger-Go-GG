package engine

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps         = 60.0
	framePeriod = 1.0 / fps * 1000.0

	// Epsilon = 0.00001
	Epsilon = 0.00001 // ~32 bits

	// DegreeToRadians converts to radians, for example, 45.0 * DegreeToRadians = radians
	DegreeToRadians = math.Pi / 180.0
)

// Vector3 contains base components
type Vector3 struct {
	X, Y, Z float64
}

// Matrix4 represents a column major opengl array.
type Matrix4 struct {
	e [16]float64

	// Rotation is in radians
	Rotation float64
	Scale    Vector3
}

// Game should be implemented by the developer
type Game interface {
	Update(float64, []uint8)
	Render(*image.RGBA)
}

// AffinePool is a pool of transforms
var AffinePool = NewAffineTransformPool(100)

// VectorsPool is a pool of transforms
var VectorsPool = NewVectorPool(100)

// Engine is the GUI and shell
type Engine struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	texture  *sdl.Texture

	Width  int32
	Height int32

	ClearColor color.RGBA

	game Game

	// Scene graph root which is always a GroupNode
	root IGroupNode

	// drawing buffer
	pixels *image.RGBA
	bounds image.Rectangle

	context *RenderContext

	// mouse
	mx int32
	my int32

	running bool

	opened bool

	nFont        *Font
	txtSimStatus *Text
	txtFPSLabel  *Text
	txtLoopLabel *Text
	txtMousePos  *Text
	dynaTxt      *DynaText
}

// NewEngine creates a new engine and initializes it.
func NewEngine(width, height int32) *Engine {
	v := new(Engine)
	v.Width = width
	v.Height = height
	v.opened = false
	v.ClearColor = color.RGBA{127, 127, 127, 255}

	v.root = NewGroupNode(nil, false)
	v.root.SetName("Root")

	return v
}

func (v *Engine) Initialize(title string) {
	v.initialize(title)

	v.opened = true
}

func (v *Engine) SetRoot(n *GroupNode) {
	v.root = n
}

func (v *Engine) GetRoot() IGroupNode {
	return v.root
}

// Start shows the display and begins event polling
func (v *Engine) Start(game Game) {
	v.game = game

	v.Run()
}

// SetFont sets the font based on path and size.
func (v *Engine) SetFont(fontPath string, size int) error {
	var err error
	v.nFont, err = NewFont(fontPath, size)
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (v *Engine) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		v.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		v.mx = t.X
		v.my = t.Y
		// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		return false // We handled it. Don't allow it to be added to the queue.
		// case *sdl.MouseButtonEvent:
		// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		// case *sdl.MouseWheelEvent:
		// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				v.running = false
			}
		}
		// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
		// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		return false
	}

	return true
}

// Run starts the polling event loop. This must run on
// the main thread.
func (v *Engine) Run() {
	v.running = true
	var frameStart time.Time
	var elapsedTime float64
	var loopTime float64

	sleepDelay := 0.0

	// Get a reference to SDL's internal keyboard state. It is updated
	// during sdl.PollEvent()
	keyState := sdl.GetKeyboardState()

	sdl.SetEventFilterFunc(v.filterEvent, nil)

	for v.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		dt := elapsedTime / 1000.0

		// Update the scene graph
		v.root.Update(dt)

		// Notify external clients of an update, perhaps for key events
		v.game.Update(dt, keyState)

		v.clearDisplay()

		// Render scene graph
		v.root.Render(v.context)

		// Notify external clients for any additional rendering
		v.game.Render(v.pixels)

		v.renderRawOverlay(elapsedTime, loopTime)

		v.renderer.Present()

		loopTime = float64(time.Since(frameStart).Nanoseconds() / 1000000.0)

		// Lock frame rate
		sleepDelay = math.Floor(framePeriod - loopTime)
		if sleepDelay > 0 {
			// fmt.Printf("%3.5f ,%3.5f, %3.5f, %3.5f \n", framePeriod, elapsedTime, sleepDelay, loopTime)
			sdl.Delay(uint32(sleepDelay))
			elapsedTime = framePeriod
		} else {
			elapsedTime = framePeriod
		}
	}
}

func (v *Engine) renderRawOverlay(elapsedTime, loopTime float64) {
	// v.texture.Update(nil, v.pixels, v.pixelPitch)
	// This takes on average 5-7ms
	v.texture.Update(nil, v.pixels.Pix, v.pixels.Stride)
	v.renderer.Copy(v.texture, nil, nil)

	v.txtFPSLabel.DrawAt(10, 10)
	f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
	v.dynaTxt.DrawAt(v.txtFPSLabel.Bounds.W+10, 10, f)

	// v.mx, v.my, _ = sdl.GetMouseState()
	v.txtMousePos.DrawAt(10, 25)
	f = fmt.Sprintf("<%d, %d>", v.mx, v.my)
	v.dynaTxt.DrawAt(v.txtMousePos.Bounds.W+10, 25, f)

	v.txtLoopLabel.DrawAt(10, 40)
	f = fmt.Sprintf("%2.2f", loopTime)
	v.dynaTxt.DrawAt(v.txtLoopLabel.Bounds.W+10, 40, f)
}

// Quit stops the engine from running, effectively shutting it down.
func (v *Engine) Quit() {
	v.running = false
}

// Close releases resources and shutsdown the engine.
// Be sure to setup a "defer x.Close()"
func (v *Engine) Close() {
	if !v.opened {
		return
	}
	var err error

	v.nFont.Destroy()
	v.txtFPSLabel.Destroy()
	v.txtMousePos.Destroy()
	v.dynaTxt.Destroy()

	log.Println("Destroying texture")
	err = v.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	v.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying window")
	err = v.window.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (v *Engine) initialize(title string) {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	v.window, err = sdl.CreateWindow(title, 100, 100,
		v.Width, v.Height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	// Using GetSurface requires using window.UpdateSurface() rather than renderer.Present.
	// v.surface, err = v.window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// v.renderer, err = sdl.CreateSoftwareRenderer(v.surface)
	// OR create renderer manually
	v.renderer, err = sdl.CreateRenderer(v.window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}

	v.texture, err = v.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, v.Width, v.Height)
	if err != nil {
		panic(err)
	}

	v.bounds = image.Rect(0, 0, int(v.Width), int(v.Height))
	v.pixels = image.NewRGBA(v.bounds)

	v.context = NewRenderContext(v.pixels)
}

// Configure view with draw objects
func (v *Engine) Configure() {
	// rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	// v.renderer.SetDrawColor(255, 127, 0, 255)
	// v.renderer.FillRect(&rect)

	v.txtSimStatus = NewText(v.nFont, v.renderer)
	err := v.txtSimStatus.SetText("Sim Status: ", sdl.Color{R: 0, G: 0, B: 255, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtFPSLabel = NewText(v.nFont, v.renderer)
	err = v.txtFPSLabel.SetText("FPS: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtMousePos = NewText(v.nFont, v.renderer)
	err = v.txtMousePos.SetText("Mouse: ", sdl.Color{R: 255, G: 127, B: 0, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtLoopLabel = NewText(v.nFont, v.renderer)
	err = v.txtLoopLabel.SetText("Loop: ", sdl.Color{R: 255, G: 127, B: 0, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.dynaTxt = NewDynaText(v.nFont, v.renderer, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}

func (v *Engine) clearDisplay() {
	for y := 0; y < int(v.Height); y++ {
		for x := 0; x < int(v.Width); x++ {
			v.pixels.SetRGBA(x, y, v.ClearColor)
		}
	}
	// v.renderer.SetDrawColor(127, 127, 127, 255)
	// v.renderer.Clear()
	// v.renderer.Present()
	// v.window.UpdateSurface()
}
