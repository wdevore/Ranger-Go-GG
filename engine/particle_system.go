package engine

import (
	"image"
)

// RenderParticle implemented by developer to render a single particle.
type RenderParticle func(particle *Particle, pixels *image.RGBA)

// TriggerParticle activates a particle defined by developer.
type TriggerParticle func(particle *Particle, system *ParticleSystem)

// ParticleSystem is a bunch of particles.
type ParticleSystem struct {
	particles []*Particle

	renderer  RenderParticle
	triggerer TriggerParticle
}

// NewParticleSystem creates a particle system
func NewParticleSystem(count int, renderer RenderParticle, triggerer TriggerParticle) *ParticleSystem {
	ps := new(ParticleSystem)
	ps.renderer = renderer
	ps.triggerer = triggerer
	ps.Initialize(count)
	return ps
}

// Initialize creates particles
func (ps *ParticleSystem) Initialize(count int) {
	ps.particles = make([]*Particle, count)
	for i := range ps.particles {
		ps.particles[i] = NewParticle()
	}
}

// Update processes all particles.
func (ps *ParticleSystem) Update(dt float64) {
	for _, p := range ps.particles {
		if p.IsAlive {
			p.Update(dt)
		}
	}
}

// Render draws a particle using the "renderer" callback.
func (ps *ParticleSystem) Render(pixels *image.RGBA) {
	for _, p := range ps.particles {
		if p.IsAlive {
			ps.renderer(p, pixels)
		}
	}
}

// TriggerParticle activates a single particle
func (ps *ParticleSystem) TriggerParticle() {
	for _, p := range ps.particles {
		if !p.IsAlive {
			p.IsAlive = true
			p.RenderColor = p.StartColor
			p.acculTime = 0
			ps.triggerer(p, ps)
			return
		}
	}
}
