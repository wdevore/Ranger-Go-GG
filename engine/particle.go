package engine

import (
	"image/color"
)

// Particle is a single particle
type Particle struct {
	Position *Vector3
	Velocity *Vector3

	// How long a particle lives after activation.
	Duration  float64
	acculTime float64

	// Colors are lerped
	StartColor, EndColor color.RGBA
	RenderColor          color.RGBA

	IsAlive bool
}

// NewParticle creates a new particle at 0,0 and white.
func NewParticle() *Particle {
	p := new(Particle)
	p.Duration = 1.0

	p.Position = NewVector3()
	p.Velocity = NewVector3()
	return p
}

// Update modifies a particle's particles.
func (ps *Particle) Update(dt float64) {
	if !ps.IsAlive {
		return
	}

	re := LinearEasing(ps.acculTime, float64(ps.StartColor.R), -(float64(ps.StartColor.R) - float64(ps.EndColor.R)), ps.Duration)
	ps.RenderColor.R = uint8(re)

	gr := LinearEasing(ps.acculTime, float64(ps.StartColor.G), -(float64(ps.StartColor.G) - float64(ps.EndColor.G)), ps.Duration)
	ps.RenderColor.G = uint8(gr)

	bl := LinearEasing(ps.acculTime, float64(ps.StartColor.B), -(float64(ps.StartColor.B) - float64(ps.EndColor.B)), ps.Duration)
	ps.RenderColor.B = uint8(bl)

	ps.RenderColor.A = ps.StartColor.A

	ps.acculTime += float64(dt)
	// fmt.Printf("%f, %f, %f: %f\n", re, gr, bl, ps.acculTime)

	ps.Position.Add(ps.Velocity)

	if ps.acculTime >= ps.Duration {
		// Particle expired
		// fmt.Printf("particle expired: %f\n", ps.acculTime)
		ps.IsAlive = false
	}
}
