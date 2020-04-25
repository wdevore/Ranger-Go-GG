package main

import (
	"image"
	"image/color"
	"math"
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

	gEngine.Initialize("Particles")

	gEngine.SetFont("../galacticstormexpand.ttf", 16)
	gEngine.Configure()

	game := newGame()
	gEngine.Start(game)
}

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

type particlesGame struct {
	particles *engine.ParticleSystem

	tKey *keyState
}

func newGame() *particlesGame {
	g := new(particlesGame)

	g.particles = engine.NewParticleSystem(200, g.particleRender, g.particleTrigger)

	g.tKey = new(keyState)
	return g
}

func (pg *particlesGame) Update(dt float64, keyState []uint8) {
	// fmt.Printf("keys: %v\n", keyState)
	if pg.tKey.isKeyUp(keyState[sdl.SCANCODE_T]) {
		pg.particles.TriggerParticle()
	}

	pg.particles.TriggerParticle()

	pg.particles.Update(dt)
}

func (pg *particlesGame) Render(pixels *image.RGBA) {
	pg.particles.Render(pixels)
}

func (pg *particlesGame) particleRender(particle *engine.Particle, pixels *image.RGBA) {
	// pixels.SetRGBA(int(particle.X-1.0), int(particle.Y), particle.RenderColor)
	// pixels.SetRGBA(int(particle.X+1.0), int(particle.Y), particle.RenderColor)
	// pixels.SetRGBA(int(particle.X), int(particle.Y-1.0), particle.RenderColor)
	// pixels.SetRGBA(int(particle.X), int(particle.Y+1.0), particle.RenderColor)
	engine.DrawRect(particle.Position.X, particle.Position.Y, 10, 10, true, particle.RenderColor, pixels)
}

func (pg *particlesGame) particleTrigger(particle *engine.Particle, system *engine.ParticleSystem) {
	// Generate a velocity vector
	v := ran.Float64()*3.0 + 1.0
	vx := math.Sin(ran.Float64() * math.Pi * 2)
	vy := math.Cos(ran.Float64() * math.Pi * 2)
	particle.Velocity.Set2Components(vx, vy)
	particle.Velocity.ScaleBy(v)

	c := uint8(ran.Float64()*64 + 32)
	particle.StartColor = color.RGBA{c, 32, 32, 255}
	r := uint8(ran.Float64()*127 + 127)
	g := uint8(ran.Float64()*127 + 127)
	b := uint8(ran.Float64()*127 + 127)
	particle.EndColor = color.RGBA{r, g, b, 255}

	particle.Position.Set2Components(float64(gEngine.Width/2), float64(gEngine.Height/2))
	particle.Duration = 2.0
}
