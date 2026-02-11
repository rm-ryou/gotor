package explorer

type Tree struct {
	node         *Node
	selectedPath string

	flatCache []*Node
	dirty     bool
}

func New(n *Node) *Tree {
	return &Tree{
		node:  n,
		dirty: true,
	}
}

func (t *Tree) Root() *Node {
	return t.node
}

func (t *Tree) SelectedPath() string {
	return t.selectedPath
}

func (t *Tree) VisibleNodes() []*Node {
	if t.dirty {
		t.rebuild()
		t.dirty = false
	}
	return t.flatCache
}

func (t *Tree) rebuild() {
	t.flatCache = t.flatCache[:0]
	if t.node.Expanded {
		for _, c := range t.node.Children {
			c.Flatten(&t.flatCache)
		}
	}
}

func (t *Tree) Toggle(node *Node, children []*Node) {
	node.Expanded = !node.Expanded

	if node.Expanded && len(node.Children) == 0 && len(children) > 0 {
		node.Children = children
	}

	if !node.Expanded && node.ContainsPath(t.selectedPath) {
		t.selectedPath = ""
	}

	t.dirty = true
}
