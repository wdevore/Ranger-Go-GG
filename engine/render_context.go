package engine

import (
	"image"
	"image/color"
)

const (
	MaxTranformedVertices = 100
)

// RenderContext is a rendering context
type RenderContext struct {
	tPoints []*Vector3

	dc *gg.Context
	// Current context
	context      *AffineTransform
	contextState *Stack
}

func NewRenderContext(image *image.RGBA) *RenderContext {
	c := new(RenderContext)
	c.contextState = NewStack(100)
	c.dc = gg.NewContextForRGBA(image)

	c.context = NewAffineTransform()
	c.tPoints = make([]*Vector3, MaxTranformedVertices)

	for i := range c.tPoints {
		c.tPoints[i] = NewVector3()
	}

	return c
}

func (c *RenderContext) Set(at *AffineTransform) {
	c.context = at
}

func (c *RenderContext) TransformContext() *AffineTransform {
	return c.context
}

func (c *RenderContext) RenderContext() *gg.Context {
	return c.dc
}

func (c *RenderContext) Save() {
	c.dc.Push()

	t := AffinePool.Pop()
	t.SetWithAT(c.context) // Copy current context and push
	c.contextState.Push(t)
}

func (c *RenderContext) Transform(at *AffineTransform) {
	AffineTransformMultiplyTo(at, c.context)
}

func (c *RenderContext) Restore() {
	t := c.contextState.Pop().(*AffineTransform)
	c.context.SetWithAT(t) // Copy to current context
	AffinePool.Push(t)

	c.dc.Pop()
}

func (c *RenderContext) DrawPolygon(vertices []*Vector3, color color.RGBA) {
	// Transform geometry for rendering
	for i, p := range vertices {
		c.context.ApplyTo(p, c.tPoints[i])
	}

	c.dc.SetColor(color)
	c.dc.MoveTo(c.tPoints[0].X, c.tPoints[0].Y)

	for i := 1; i < len(vertices); i++ {
		c.dc.LineTo(c.tPoints[i].X, c.tPoints[i].Y)
	}
	// for _, t := range c.tPoints[1:] {
	// }

	c.dc.ClosePath()
	c.dc.Fill()
}
