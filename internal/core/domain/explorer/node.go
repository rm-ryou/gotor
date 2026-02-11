package explorer

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
