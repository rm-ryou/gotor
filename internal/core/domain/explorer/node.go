package explorer

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Depth    int
	Expanded bool
	Children []*Node
}
