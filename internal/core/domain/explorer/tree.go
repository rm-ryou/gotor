package explorer

import "path/filepath"

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

func (t *Tree) Select(node *Node) {
	if node.IsDir {
		return
	}
	t.selectedPath = node.Path
}

func (t *Tree) ClearSelection() {
	t.selectedPath = ""
}

func (t *Tree) Reset(node *Node) {
	t.node = node
	t.selectedPath = ""
	t.flatCache = t.flatCache[:0]
	t.dirty = true
}

func (t *Tree) FindNode(path string) *Node {
	return findNode(t.node, path)
}

func findNode(node *Node, path string) *Node {
	if node.Path == path {
		return node
	}

	rel, err := filepath.Rel(node.Path, path)
	if err != nil || len(rel) >= 2 && rel[:2] == ".." {
		return nil
	}

	for _, c := range node.Children {
		if found := findNode(c, path); found != nil {
			return found
		}
	}
	return nil
}
