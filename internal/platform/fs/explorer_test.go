package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "internal"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "cmd"), 0o755); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"main.go", "go.mod", ".gitignore"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o6444); err != nil {
			t.Fatal(err)
		}
	}

	return root
}

func Test_ReadDir_HiddenExcluded(t *testing.T) {
	t.Parallel()
	root := setupTestDir(t)

	reader := New(false)
	nodes, err := reader.ReadDir(root, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, n := range nodes {
		if n.Name == ".gitignore" {
			t.Errorf("hidden file should be excluded when showHidden = false")
		}
	}
}

func Test_ReadDir_HiddenIncluded(t *testing.T) {
	t.Parallel()
	root := setupTestDir(t)

	reader := New(true)
	nodes, err := reader.ReadDir(root, 1)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}

	found := false
	for _, n := range nodes {
		if n.Name == ".gitignore" {
			found = true
		}
	}
	if !found {
		t.Error("hidden file should be included when showHidden = true")
	}
}

func Test_ReadDir_SortOrder(t *testing.T) {
	t.Parallel()
	root := setupTestDir(t)

	reader := New(false)
	nodes, err := reader.ReadDir(root, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, n := range nodes {
		if !n.IsDir {
			for _, rest := range nodes[i:] {
				if rest.IsDir {
					t.Errorf("dirs must come first. but %s appears after file %s", rest.Name, n.Name)
				}
			}
			break
		}
	}

	var dirs []string
	for _, n := range nodes {
		if n.IsDir {
			dirs = append(dirs, n.Name)
		}
	}
	for i := 1; i < len(dirs); i++ {
		if dirs[i-1] > dirs[i] {
			t.Errorf("not sorted: %q > %q", dirs[i-1], dirs[i])
		}
	}
}

func Test_ReadDir_NotExist(t *testing.T) {
	t.Parallel()

	reader := New(false)
	_, err := reader.ReadDir("/path/that/does/not/exist", 1)
	if err == nil {
		t.Error("expected error for non-existent path, act nil")
	}
}
