package engine

type IGroupNode interface {
	INode
	Add(n INode) // Last node added is render underneath
	Remove(n INode)
	Find(n INode) (f int, fno INode)
}

// GroupNode is a collection of nodes
type GroupNode struct {
	BaseNode // is-a

	nodes []INode
}

func NewGroupNode(parent IGroupNode, autoAdd bool) IGroupNode {
	g := new(GroupNode)
	g.Initialize()
	g.parent = parent
	g.nodes = []INode{}

	if autoAdd {
		parent.Add(g)
	}
	return g
}

func (gn *GroupNode) Add(n INode) {
	gn.nodes = append(gn.nodes, n)
}

func (gn *GroupNode) Remove(n INode) {
	j, no := gn.Find(n)
	if no != nil {
		//gn.nodes = append(gn.nodes[:j], gn.nodes[j+1:]...)	// <-- this may leak
		// https://github.com/golang/go/wiki/SliceTricks
		copy(gn.nodes[j:], gn.nodes[j+1:])
		gn.nodes[len(gn.nodes)-1] = nil // or the zero value of T
		gn.nodes = gn.nodes[:len(gn.nodes)-1]
	}
}

func (gn *GroupNode) Find(n INode) (f int, fno INode) {
	for i, no := range gn.nodes {
		if no == n {
			f = i
			fno = no
			break
		}
	}
	return f, fno
}

func (gn *GroupNode) Update(dt float64) {
	// Update properties of the group node
	gn.BaseNode.Update(dt)

	for _, n := range gn.nodes {
		n.Update(dt)
	}

	// Update node's transform if dirty
}

func (gn *GroupNode) Render(context *RenderContext) {
	if !gn.IsVisible() {
		return
	}

	// Save context state first
	context.Save()

	// Append this node's transform onto the context and then render
	context.Transform(gn.calcTransform())

	for _, n := range gn.nodes {
		// fmt.Printf("GroupNode render: %s\n", n)
		n.Render(context)
	}

	// Now draw this node if it has an geometry, typically it doesn't
	gn.Draw(context)

	context.Restore()
}

func (gn *GroupNode) Draw(context *RenderContext) {
	// gn.BaseNode.Draw(context)
	// fmt.Println("GroupNode::Draw")

	// context.DrawPolygon(0, 0, 1, 1, true, n.SolidColor)
	// Draw using transformed geometry
}
