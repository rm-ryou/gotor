package explorer

import "path/filepath"

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Depth    int
	Expanded bool
	Children []*Node
}

func (n *Node) Flatten(out *[]*Node) {
	*out = append(*out, n)
	if n.IsDir && n.Expanded {
		for _, c := range n.Children {
			c.Flatten(out)
		}
	}
}

func (n *Node) ContainsPath(path string) bool {
	rel, err := filepath.Rel(n.Path, path)
	if err != nil {
		return false
	}

	if rel == "." || len(rel) >= 2 && rel[:2] == ".." {
		return false
	}
	return true
}
