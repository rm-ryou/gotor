package usecase

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/rm-ryou/gotor/internal/core/domain/cursor"
	"github.com/rm-ryou/gotor/internal/core/domain/document"
)

type FileIO interface {
	document.Reader
	document.Writer
	document.Deleter
}

type Editor struct {
	fileIO FileIO

	doc      document.Document
	cursor   *cursor.Cursor
	filePath string
	dirty    bool
}

func NewEditor(fio FileIO) *Editor {
	return &Editor{
		fileIO:   fio,
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
	content, err := e.fileIO.Read(path)
	if err != nil {
		return NewError("Failed to open the file.", err)
	}

	e.doc = document.NewFromText(content)
	e.cursor = cursor.New(0, 0)
	e.filePath = path
	e.dirty = false

	return nil
}

func (e *Editor) InsertText(text string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if text == "" {
		return
	}

	offset := e.cursorOffset()
	e.doc.Insert(offset, text)
	e.setCursorFromOffset(offset + len(text))
	e.dirty = true
}

func (e *Editor) DeleteBackward() bool {
	offset := e.cursorOffset()
	if offset == 0 {
		return false
	}

	start := previousRuneStart(e.doc.Text(), offset)
	e.doc.Delete(start, offset)
	e.setCursorFromOffset(start)
	e.dirty = true
	return true
}

func (e *Editor) Save() error {
	if e.filePath == "" {
		return NewError("No file is selected for saving.", errors.New("no file path set"))
	}
	if err := e.fileIO.Write(e.filePath, e.doc.Text()); err != nil {
		return NewError("Failed to save the file.", err)
	}
	e.dirty = false
	return nil
}

func (e *Editor) SaveAs(path string) error {
	if err := e.fileIO.Write(path, e.doc.Text()); err != nil {
		return NewError("Failed to save the file.", err)
	}
	e.filePath = path
	e.dirty = false
	return nil
}

func (e *Editor) DeleteFile() error {
	if e.filePath == "" {
		return NewError("No file is selected for deletion.", errors.New("no file path set"))
	}
	if err := e.fileIO.Delete(e.filePath); err != nil {
		return NewError("Failed to delete the file.", err)
	}
	e.NewFile()
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

func (e *Editor) cursorOffset() int {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return 0
	}

	row := e.cursor.Row
	if row < 0 {
		row = 0
	}
	if row >= len(lines) {
		row = len(lines) - 1
	}

	offset := 0
	for i := 0; i < row; i++ {
		offset += len(lines[i]) + 1
	}

	return offset + byteOffsetForColumn(lines[row], e.cursor.Col)
}

func (e *Editor) setCursorFromOffset(offset int) {
	if offset < 0 {
		offset = 0
	}

	text := e.doc.Text()
	if offset > len(text) {
		offset = len(text)
	}

	row := 0
	col := 0
	read := 0

	for _, r := range text {
		if read >= offset {
			break
		}

		if r == '\n' {
			row++
			col = 0
		} else {
			col++
		}

		read += utf8.RuneLen(r)
	}

	e.cursor.MoveTo(row, col)
	e.clampCursor()
}

func byteOffsetForColumn(s string, col int) int {
	if col <= 0 {
		return 0
	}

	offset := 0
	count := 0
	for _, r := range s {
		if count >= col {
			break
		}
		offset += utf8.RuneLen(r)
		count++
	}

	return offset
}

func previousRuneStart(s string, offset int) int {
	if offset <= 0 {
		return 0
	}
	if offset > len(s) {
		offset = len(s)
	}

	start := 0
	for i := range s {
		if i >= offset {
			break
		}
		start = i
	}

	return start
}
