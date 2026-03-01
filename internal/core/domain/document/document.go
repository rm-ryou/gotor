package document

type Document interface {
	Lines() []string
	LineCount() int
	Line(row int) string
	Text() string
	Insert(offset int, text string)
	Delete(start, end int)
}

type Reader interface {
	Read(path string) (string, error)
}

type Writer interface {
	Write(path, content string) error
}

type Deleter interface {
	Delete(path string) error
}

func New() Document {
	return newPieceTable()
}

func NewFromText(text string) Document {
	return newPieceTableFromText(text)
}
