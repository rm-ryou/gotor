package document

type Document interface {
	Lines() []string
	LineCount() int
	Line(row int) string
	Text() string
}

type Reader interface {
	Read(path string) (string, error)
}

func New() Document {
	return newPieceTable()
}

func NewFromText(text string) Document {
	return newPieceTableFromText(text)
}
