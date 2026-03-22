package usecase

import "testing"

func TestEditorInsertTextAndDeleteBackward(t *testing.T) {
	t.Parallel()

	fio := &stubFileIO{content: "hello"}
	editor := NewEditor(fio)
	if err := editor.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	editor.MoveCursorRight()
	editor.MoveCursorRight()
	editor.InsertText("X")

	if got := editor.Document().Text(); got != "heXllo" {
		t.Fatalf("Text() after insert = %q, want %q", got, "heXllo")
	}
	if got := editor.Cursor().Col; got != 3 {
		t.Fatalf("cursor col after insert = %d, want 3", got)
	}
	if !editor.IsDirty() {
		t.Fatal("IsDirty() = false, want true after insert")
	}

	if !editor.DeleteBackward() {
		t.Fatal("DeleteBackward() = false, want true")
	}
	if got := editor.Document().Text(); got != "hello" {
		t.Fatalf("Text() after delete = %q, want %q", got, "hello")
	}
	if got := editor.Cursor().Col; got != 2 {
		t.Fatalf("cursor col after delete = %d, want 2", got)
	}
}

func TestEditorSaveAndDeleteFile(t *testing.T) {
	t.Parallel()

	fio := &stubFileIO{content: "hello"}
	editor := NewEditor(fio)
	if err := editor.OpenFile("test.txt"); err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	editor.InsertText("!")
	if err := editor.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if got := fio.writes["test.txt"]; got != "!hello" {
		t.Fatalf("saved content = %q, want %q", got, "!hello")
	}
	if editor.IsDirty() {
		t.Fatal("IsDirty() = true, want false after save")
	}

	if err := editor.DeleteFile(); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}
	if got := fio.deletedPath; got != "test.txt" {
		t.Fatalf("deleted path = %q, want %q", got, "test.txt")
	}
	if got := editor.FilePath(); got != "" {
		t.Fatalf("FilePath() after delete = %q, want empty", got)
	}
}

func TestEditorSaveWithoutPath(t *testing.T) {
	t.Parallel()

	editor := NewEditor(&stubFileIO{})
	if err := editor.Save(); err == nil {
		t.Fatal("Save() error = nil, want error")
	}
}

type stubFileIO struct {
	content     string
	writes      map[string]string
	deletedPath string
	readErr     error
	writeErr    error
	deleteErr   error
}

func (s *stubFileIO) Read(string) (string, error) {
	if s.readErr != nil {
		return "", s.readErr
	}
	return s.content, nil
}

func (s *stubFileIO) Write(path, content string) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	if s.writes == nil {
		s.writes = map[string]string{}
	}
	s.writes[path] = content
	return nil
}

func (s *stubFileIO) Delete(path string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	s.deletedPath = path
	return nil
}

var _ FileIO = (*stubFileIO)(nil)
