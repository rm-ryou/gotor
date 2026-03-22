package features

import (
	"errors"
	"testing"

	"gioui.org/io/key"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/gio/config"
)

func TestDisplayColumnForCursor(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		limit int
		want  int
	}{
		{name: "start of line", line: "hello", limit: 0, want: 0},
		{name: "plain text", line: "hello", limit: 3, want: 3},
		{name: "tab expands to tab stop", line: "\tab", limit: 1, want: 4},
		{name: "tab after text uses next stop", line: "a\tb", limit: 2, want: 4},
		{name: "limit past line length", line: "ab", limit: 10, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayColumnForCursor(tt.line, tt.limit, config.Default().Editor.TabWidth)
			if got != tt.want {
				t.Fatalf("displayColumnForCursor(%q, %d, %d) = %d, want %d", tt.line, tt.limit, config.Default().Editor.TabWidth, got, tt.want)
			}
		})
	}
}

func TestHandleKeyEvent(t *testing.T) {
	uc := usecase.NewEditor(&stubReader{text: "ab\ncd"})
	if err := uc.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	ev := &EditorView{uc: uc}

	if !ev.handleKeyEvent(key.Event{Name: key.NameRightArrow}) {
		t.Fatal("handleKeyEvent(right) = false, want true")
	}
	if got := uc.Cursor().Col; got != 1 {
		t.Fatalf("cursor col after right = %d, want 1", got)
	}

	if !ev.handleKeyEvent(key.Event{Name: key.NameDownArrow}) {
		t.Fatal("handleKeyEvent(down) = false, want true")
	}
	if got := uc.Cursor().Row; got != 1 {
		t.Fatalf("cursor row after down = %d, want 1", got)
	}

	if !ev.handleKeyEvent(key.Event{Name: key.NameLeftArrow}) {
		t.Fatal("handleKeyEvent(left) = false, want true")
	}
	if got := uc.Cursor().Col; got != 0 {
		t.Fatalf("cursor col after left = %d, want 0", got)
	}

	if !ev.handleKeyEvent(key.Event{Name: key.NameUpArrow}) {
		t.Fatal("handleKeyEvent(up) = false, want true")
	}
	if got := uc.Cursor().Row; got != 0 {
		t.Fatalf("cursor row after up = %d, want 0", got)
	}

	if !ev.handleKeyEvent(key.Event{Name: key.NameEnter}) {
		t.Fatal("handleKeyEvent(enter) = false, want true")
	}
	if got := uc.Cursor().Row; got != 1 {
		t.Fatalf("cursor row after enter = %d, want 1", got)
	}
	if got := uc.Cursor().Col; got != 0 {
		t.Fatalf("cursor col after enter = %d, want 0", got)
	}
}

func TestHandleKeyEventSaveShortcut(t *testing.T) {
	reader := &stubReader{text: "ab"}
	uc := usecase.NewEditor(reader)
	if err := uc.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}
	uc.InsertText("X")

	ev := &EditorView{uc: uc}
	if ev.handleKeyEvent(key.Event{Name: "S", Modifiers: key.ModShortcut}) {
		t.Fatal("handleKeyEvent(save) = true, want false")
	}

	if reader.lastWritePath != "test.txt" {
		t.Fatalf("lastWritePath = %q, want %q", reader.lastWritePath, "test.txt")
	}
	if reader.lastWriteContent != "Xab" {
		t.Fatalf("lastWriteContent = %q, want %q", reader.lastWriteContent, "Xab")
	}
	if uc.IsDirty() {
		t.Fatal("IsDirty() = true, want false after save shortcut")
	}
}

func TestHandleTextInput(t *testing.T) {
	uc := usecase.NewEditor(&stubReader{text: "ab"})
	if err := uc.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	ev := &EditorView{uc: uc}

	if !ev.handleTextInput("X") {
		t.Fatal("handleTextInput(\"X\") = false, want true")
	}
	if got := uc.Document().Text(); got != "Xab" {
		t.Fatalf("Text() after input = %q, want %q", got, "Xab")
	}
	if got := uc.Cursor().Col; got != 1 {
		t.Fatalf("cursor col after input = %d, want 1", got)
	}

	if ev.handleTextInput("") {
		t.Fatal("handleTextInput(\"\") = true, want false")
	}
}

func TestHandleKeyEventSaveShortcutReportsError(t *testing.T) {
	uc := usecase.NewEditor(&stubReaderWithError{text: "ab", writeErr: errors.New("disk full")})
	if err := uc.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	var reported error
	ev := &EditorView{
		uc: uc,
		OnError: func(err error) {
			reported = err
		},
	}

	ev.handleKeyEvent(key.Event{Name: "S", Modifiers: key.ModShortcut})

	if reported == nil {
		t.Fatal("reported error = nil, want error")
	}
	if got := usecase.MessageFor(reported); got != "Failed to save the file." {
		t.Fatalf("MessageFor(reported) = %q, want %q", got, "Failed to save the file.")
	}
}

type stubReader struct {
	text             string
	lastWritePath    string
	lastWriteContent string
}

func (s stubReader) Read(string) (string, error) {
	return s.text, nil
}

func (s *stubReader) Write(path, content string) error {
	s.lastWritePath = path
	s.lastWriteContent = content
	return nil
}

func (s stubReader) Delete(string) error {
	return nil
}

type stubReaderWithError struct {
	text     string
	writeErr error
}

func (s stubReaderWithError) Read(string) (string, error) {
	return s.text, nil
}

func (s *stubReaderWithError) Write(string, string) error {
	return s.writeErr
}

func (s stubReaderWithError) Delete(string) error {
	return nil
}
