package explorer

type FSReader interface {
	ReadDir(path string, depth int) ([]*Node, error)
}

type FSWriter interface {
	CreateFile(path string) error
}
