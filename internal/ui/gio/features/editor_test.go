package features

import (
	"testing"

	"gioui.org/io/key"
	"github.com/rm-ryou/gotor/internal/core/usecase"
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
			got := displayColumnForCursor(tt.line, tt.limit, tabWidth)
			if got != tt.want {
				t.Fatalf("displayColumnForCursor(%q, %d, %d) = %d, want %d", tt.line, tt.limit, tabWidth, got, tt.want)
			}
		})
	}
}

func TestHandleKeyEvent(t *testing.T) {
	uc := usecase.NewEditor(stubReader{text: "ab\ncd"})
	if err := uc.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	ev := &EditorView{uc: uc}

	if !ev.handleKeyEvent(key.NameRightArrow) {
		t.Fatal("handleKeyEvent(right) = false, want true")
	}
	if got := uc.Cursor().Col; got != 1 {
		t.Fatalf("cursor col after right = %d, want 1", got)
	}

	if !ev.handleKeyEvent(key.NameDownArrow) {
		t.Fatal("handleKeyEvent(down) = false, want true")
	}
	if got := uc.Cursor().Row; got != 1 {
		t.Fatalf("cursor row after down = %d, want 1", got)
	}

	if !ev.handleKeyEvent(key.NameLeftArrow) {
		t.Fatal("handleKeyEvent(left) = false, want true")
	}
	if got := uc.Cursor().Col; got != 0 {
		t.Fatalf("cursor col after left = %d, want 0", got)
	}

	if !ev.handleKeyEvent(key.NameUpArrow) {
		t.Fatal("handleKeyEvent(up) = false, want true")
	}
	if got := uc.Cursor().Row; got != 0 {
		t.Fatalf("cursor row after up = %d, want 0", got)
	}

	if !ev.handleKeyEvent(key.NameEnter) {
		t.Fatal("handleKeyEvent(enter) = false, want true")
	}
	if got := uc.Cursor().Row; got != 1 {
		t.Fatalf("cursor row after enter = %d, want 1", got)
	}
	if got := uc.Cursor().Col; got != 0 {
		t.Fatalf("cursor col after enter = %d, want 0", got)
	}
}

func TestHandleTextInput(t *testing.T) {
	uc := usecase.NewEditor(stubReader{text: "ab"})
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

type stubReader struct {
	text string
}

func (s stubReader) Read(string) (string, error) {
	return s.text, nil
}

func (s stubReader) Write(string, string) error {
	return nil
}

func (s stubReader) Delete(string) error {
	return nil
}
