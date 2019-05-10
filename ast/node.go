package ast

// Node of AST
type Node struct {
	ast *AST

	op       int // Inst
	father   *Node
	children []*Node
	next     *Node
	data     Data

	Reset bool // Error Recovery
}

// Root Node Builder
func (a *AST) Root(op int) *Node {
	pr := a.instParam[op]
	var children []*Node
	if pr == nil {
		children = make([]*Node, 0)
	} else {
		children = make([]*Node, 0, pr.Max)
	}
	a.root = &Node{
		ast:      a,
		op:       op,
		father:   nil,
		children: children,

		Reset: true,
	}
	return a.root
}

// Child Node Builder
func (n *Node) Child(op int) *Node {
	r := n.ast.instParam[n.op]
	if r != nil && len(n.children) >= r.Max {
		panic(ErrOverChildrenCap)
	}

	pr := n.ast.instParam[op]
	var children []*Node
	if pr == nil || pr.Max < 0 {
		children = make([]*Node, 0)
	} else {
		children = make([]*Node, 0, pr.Max)
	}
	c := &Node{
		ast:      n.ast,
		op:       op,
		father:   n,
		children: children,
	}
	l := len(n.children)
	if l > 0 {
		n.children[l-1].next = c
	}
	n.children = append(n.children, c)
	return c
}

// Data Loader
func (n *Node) Data(data interface{}) *Node {
	n.data = data
	return n
}

// Next Node
func (n *Node) Next() *Node {
	return n.next
}
