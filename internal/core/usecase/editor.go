package usecase

import (
	"fmt"

	"github.com/rm-ryou/gotor/internal/core/domain/document"
)

type Editor struct {
	reader document.Reader

	doc      document.Document
	filePath string
	dirty    bool
}

func NewEditor(r document.Reader) *Editor {
	return &Editor{
		reader:   r,
		doc:      document.New(),
		filePath: "",
		dirty:    false,
	}
}

func (e *Editor) Document() document.Document {
	return e.doc
}

func (e *Editor) FilePath() string {
	return e.filePath
}

func (e *Editor) IsDirty() bool {
	return e.dirty
}

func (e *Editor) NewFile() {
	e.doc = document.New()
	e.filePath = ""
	e.dirty = false
}

func (e *Editor) OpenFile(path string) error {
	content, err := e.reader.Read(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}

	e.doc = document.NewFromText(content)
	e.filePath = path
	e.dirty = false

	return nil
}
