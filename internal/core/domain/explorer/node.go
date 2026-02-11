package explorer

import (
	"path/filepath"
	"strings"
)

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
	return isUnder(n.Path, path)
}

func isUnder(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}

	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	return true
}
