package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/GameEngine/engine"
)

// Note: You may need a ram drive:
// export TMPDIR="/Volumes/RAMDisk"
// diskutil erasevolume HFS+ 'RAMDisk' `hdiutil attach -nomount ram://2097152`

var gEngine *engine.Engine
var ran = rand.New(rand.NewSource(31))

func main() {
	gEngine = engine.NewEngine(512, 256)
	defer gEngine.Close()

	gEngine.Initialize("SimpleNode")

	err := gEngine.SetFont("../galacticstormexpand.ttf", 12)
	if err != nil {
		panic(err)
	}
	gEngine.Configure()

	game := newGame()
	gEngine.Start(game)
}

type aGame struct {
	tKey *keyState

	wgroup engine.IGroupNode
	white  engine.INode
	ogroup engine.IGroupNode
	orange engine.INode
	blue   engine.INode

	angle float64
}

func newGame() *aGame {
	g := new(aGame)
	g.tKey = new(keyState)

	// g.build1()
	g.build2()

	return g
}

// The orange rectangle orbits white triangle
func (pg *aGame) build2() {
	root := gEngine.GetRoot()

	pg.wgroup = engine.NewGroupNode(root, true)
	pg.wgroup.SetName("whiteGroup")
	pg.wgroup.SetPositionBy2Comp(150, 150)

	pg.white = engine.NewRectangleNode(pg.wgroup, true, true)
	pg.white.SetName("WhiteRect")
	pg.white.SetScaleUniform(25)

	pg.ogroup = engine.NewGroupNode(pg.wgroup, true)
	pg.ogroup.SetName("orangeGroup")

	pg.orange = engine.NewRectangleNode(pg.ogroup, true, true)
	pg.orange.SetName("OrangeRect")
	pg.orange.SetColor(color.RGBA{255, 127, 0, 255})
	pg.orange.SetPositionBy2Comp(50, 0)
	pg.orange.SetScaleUniform(15)

	pg.angle = 0.0
}

func (pg *aGame) build1() {
	root := gEngine.GetRoot()

	pg.wgroup = engine.NewGroupNode(root, true)
	pg.wgroup.SetName("g1")
	// root.Add(pg.wgroup)

	pg.white = engine.NewRectangleNode(pg.wgroup, true, true)
	pg.white.SetName("White")
	pg.white.SetPositionBy2Comp(100, 100)
	pg.white.SetScaleUniform(25)

	pg.orange = engine.NewRectangleNode(root, true, true)
	pg.orange.SetName("Orange")
	pg.orange.SetColor(color.RGBA{255, 127, 0, 255})
	pg.orange.SetPositionBy2Comp(150, 100)
	pg.orange.SetScaleUniform(25)

	pg.angle = 0.0

	// root.Add(g.rect2)
}

func (pg *aGame) Update(dt float64, keyState []uint8) {
	// fmt.Printf("keys: %v\n", keyState)
	if pg.tKey.isKeyUp(keyState[sdl.SCANCODE_T]) {
	}

	pg.white.SetRotationByDegree(-pg.angle / 4)
	pg.ogroup.SetRotationByDegree(pg.angle)
	pg.angle += 2.0
}

func (pg *aGame) Render(pixels *image.RGBA) {
}

// -----------------------------------------------------------------------
// Keys
// -----------------------------------------------------------------------
type keyState struct {
	prevState uint8
}

func (ks *keyState) isKeyUp(state uint8) bool {
	up := false
	if ks.prevState == 0 && state == 1 {
		up = true
	}

	ks.prevState = state
	return up
}
