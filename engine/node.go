package engine

import (
	"fmt"
	"image/color"
)

type Drawer func(*RenderContext)

// Node is the base node.
type INode interface {
	Update(dt float64)
	Render(context *RenderContext)

	Position() *Vector3
	SetPosition(*Vector3)
	SetPositionBy2Comp(x, y float64)

	Scale() *Vector3
	SetScale(*Vector3)
	SetScaleUniform(float64)

	Rotation() float64
	SetRotation(float64)
	SetRotationByDegree(float64)

	SetInvisible()
	SetVisible()
	IsVisible() bool

	SetColor(color.RGBA)
	Name() string
	SetName(string)

	calcTransform() *AffineTransform

	String() string
}

type BaseNode struct {
	name      string
	parent    IGroupNode
	transform *AffineTransform
	dirty     bool
	visible   bool

	position *Vector3
	scale    *Vector3
	rotation float64

	SolidColor color.RGBA

	drawer Drawer
}

func (n *BaseNode) Initialize() {
	n.dirty = true
	n.visible = true
	n.position = NewVector3()

	n.scale = NewVector3()
	n.scale.Set2Components(1.0, 1.0)

	n.SolidColor = color.RGBA{255, 255, 255, 255}
	n.transform = NewAffineTransform()
}

func (n *BaseNode) SetColor(color color.RGBA) {
	n.SolidColor = color
}

func (n *BaseNode) Position() *Vector3 {
	return n.position
}

func (n *BaseNode) SetPosition(v *Vector3) {
	n.dirty = true
	n.position.Set2Components(v.X, v.Y)
}

func (n *BaseNode) SetPositionBy2Comp(x, y float64) {
	n.dirty = true
	n.position.Set2Components(x, y)
}

func (n *BaseNode) Scale() *Vector3 {
	return n.position
}

func (n *BaseNode) SetScale(v *Vector3) {
	n.dirty = true
	n.scale.Set2Components(v.X, v.Y)
}

func (n *BaseNode) SetScaleUniform(s float64) {
	n.dirty = true
	n.scale.ScaleBy(s)
}

func (n *BaseNode) Rotation() float64 {
	return n.rotation
}

// +angle yields CW rotation
func (n *BaseNode) SetRotation(angle float64) {
	n.dirty = true
	n.rotation = angle
}

// +angle yields CW rotation
func (n *BaseNode) SetRotationByDegree(angle float64) {
	n.dirty = true
	n.rotation = angle * DegreeToRadians
}

func (n *BaseNode) Name() string {
	return n.name
}

func (n *BaseNode) SetName(s string) {
	n.name = s
}

func (n *BaseNode) SetVisible() {
	n.visible = true
}

func (n *BaseNode) SetInvisible() {
	n.visible = false
}

func (n *BaseNode) IsVisible() bool {
	return n.visible
}

// Update node
func (n *BaseNode) Update(dt float64) {
	// fmt.Println("Node::Update")
	// Update properties of the node

	// Update node's transform if dirty
}

// Render node
func (n *BaseNode) Render(context *RenderContext) {
	if !n.visible {
		return
	}

	// Save context state first
	context.Save()

	// if n.parent != nil {
	// 	context.SetWithAT(n.parent.transform)
	// } else {
	// 	context.ToIdentity()
	// }

	// Append this node's transform onto the context and then render
	context.Transform(n.calcTransform())

	// n.Draw(context)
	n.drawer(context)

	// Restores
	context.Restore()
}

/**
 * Returns a matrix that represents this [Node]'s local-space
 * transform.
 */
func (n *BaseNode) calcTransform() *AffineTransform {
	// Note: We could check each behaviors for dirty but that would be
	// expensive especially in this core method. So instead each
	// behavior is tightly bound to this BaseNode, and the behavior updates
	// the "dirty" state. The down side is that each behavior needs a
	// reference to a BaseNode that it affects.
	if n.dirty {
		n.transform.ToIdentity()

		n.transform.Translate(n.position.X, n.position.Y)

		if n.rotation != 0.0 {
			n.transform.Rotate(n.rotation)
		}

		if n.scale.X != 1.0 || n.scale.Y != 1.0 {
			n.transform.Scale(n.scale.X, n.scale.Y)
		}

		//print("BaseNode.calcTransform\n ${transform}, tag:$tag");
		n.dirty = false
	}

	return n.transform
}

func (n BaseNode) String() string {
	return fmt.Sprintf("'%s': %v", n.name, n.position)
}

// -----------------------------------------------------------------
// Basic preconfigured nodes
// -----------------------------------------------------------------

// -----------------------------------------------------------------
// Rectangle
// -----------------------------------------------------------------

type RectangleNode struct {
	BaseNode // is-a

	centered bool
	vertices []*Vector3
}

func NewRectangleNode(parent IGroupNode, centered, autoAdd bool) INode {
	g := new(RectangleNode)
	g.Initialize()
	g.centered = centered
	g.parent = parent

	if autoAdd {
		g.parent.Add(g)
	}

	g.vertices = make([]*Vector3, 4)
	g.vertices[0] = NewVector3()
	g.vertices[1] = NewVector3()
	g.vertices[2] = NewVector3()
	g.vertices[3] = NewVector3()

	if centered {
		g.vertices[0].Set2Components(-0.5, -0.5)
		g.vertices[1].Set2Components(-0.5, 0.5)
		g.vertices[2].Set2Components(0.5, 0.5)
		g.vertices[3].Set2Components(0.5, -0.5)
	} else {
		g.vertices[0].Set2Components(0.0, 0.0)
		g.vertices[1].Set2Components(0.0, 1.0)
		g.vertices[2].Set2Components(1.0, 1.0)
		g.vertices[3].Set2Components(1.0, 0.0)
	}

	g.drawer = g.Draw

	return g
}

func (n *RectangleNode) Update(dt float64) {
}

func (n *RectangleNode) Render(context *RenderContext) {
	if !n.IsVisible() {
		return
	}
	n.BaseNode.Render(context)

	// n.Draw(context)
}

func (n *RectangleNode) Draw(context *RenderContext) {
	context.DrawPolygon(n.vertices, n.SolidColor)
}
