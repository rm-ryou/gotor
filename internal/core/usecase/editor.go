package usecase

import (
	"fmt"

	"github.com/rm-ryou/gotor/internal/core/domain/cursor"
	"github.com/rm-ryou/gotor/internal/core/domain/document"
)

type Editor struct {
	reader document.Reader

	doc      document.Document
	cursor   *cursor.Cursor
	filePath string
	dirty    bool
}

func NewEditor(r document.Reader) *Editor {
	return &Editor{
		reader:   r,
		doc:      document.New(),
		cursor:   cursor.New(0, 0),
		filePath: "",
		dirty:    false,
	}
}

func (e *Editor) Document() document.Document {
	return e.doc
}

func (e *Editor) Cursor() *cursor.Cursor {
	return e.cursor
}

func (e *Editor) FilePath() string {
	return e.filePath
}

func (e *Editor) IsDirty() bool {
	return e.dirty
}

func (e *Editor) NewFile() {
	e.doc = document.New()
	e.cursor = cursor.New(0, 0)
	e.filePath = ""
	e.dirty = false
}

func (e *Editor) OpenFile(path string) error {
	content, err := e.reader.Read(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}

	e.doc = document.NewFromText(content)
	e.cursor = cursor.New(0, 0)
	e.filePath = path
	e.dirty = false

	return nil
}

func (e *Editor) clampCursor() {
	if e.cursor.Row < 0 {
		e.cursor.Row = 0
	}
	if e.cursor.Row >= e.doc.LineCount() {
		e.cursor.Row = e.doc.LineCount() - 1
	}

	lineLen := len([]rune(e.doc.Line(e.cursor.Row)))
	if e.cursor.Col < 0 {
		e.cursor.Col = 0
	}
	if e.cursor.Col > lineLen {
		e.cursor.Col = lineLen
	}
}

func (e *Editor) MoveCursorUp() {
	e.cursor.MoveUp()
	e.clampCursor()
}

func (e *Editor) MoveCursorDown() {
	e.cursor.MoveDown()
	e.clampCursor()
}

func (e *Editor) MoveCursorLeft() {
	e.cursor.MoveLeft()
	e.clampCursor()
}

func (e *Editor) MoveCursorRight() {
	e.cursor.MoveRight()
	e.clampCursor()
}

func (e *Editor) MoveCursorToLineStart() {
	e.cursor.MoveToStartLine()
	e.clampCursor()
}

func (e *Editor) MoveCursorToLineEnd() {
	lineLen := len([]rune(e.doc.Line(e.cursor.Row)))
	e.cursor.MoveToEndLine(lineLen)
	e.clampCursor()
}
