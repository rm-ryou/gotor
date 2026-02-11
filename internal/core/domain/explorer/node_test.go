package explorer

import "testing"

func Test_Flatten_FileNode(t *testing.T) {
	t.Parallel()

	var act []*Node
	n := &Node{
		Name:  "main.go",
		Path:  "/path/main.go",
		IsDir: false,
	}

	n.Flatten(&act)

	if len(act) != 1 {
		t.Fatalf("want: 1, act: %d", len(act))
	}
	if act[0].Name != n.Name {
		t.Errorf("want: %s, act: %s", n.Name, act[0].Name)
	}
}

func Test_Flatten_CollapsedDir(t *testing.T) {
	t.Parallel()

	var act []*Node
	dir := &Node{
		Name:     "internal",
		Path:     "/path/internal",
		IsDir:    true,
		Expanded: false,
		Children: []*Node{
			{Name: "hoge.go", Path: "/path/internal/hoge.go"},
		},
	}

	dir.Flatten(&act)

	if len(act) != 1 {
		t.Fatalf("want: 1, act: %d", len(act))
	}
	if act[0].Name != dir.Name {
		t.Errorf("want: %s, act: %s", dir.Name, act[0].Name)
	}
}

func Test_Flatten_ExpandedDir(t *testing.T) {
	t.Parallel()

	var act []*Node
	dir := &Node{
		Name:     "internal",
		Path:     "path/internal",
		IsDir:    true,
		Expanded: true,
		Children: []*Node{
			{
				Name:     "core",
				Path:     "/path/internal/core",
				IsDir:    true,
				Expanded: true,
				Children: []*Node{
					{Name: "domain.go", Path: "/path/internal/core/domain.go"},
				},
			},
			{Name: "hoge.go", Path: "/path/internal/hoge.go"},
		},
	}

	dir.Flatten(&act)

	wantNames := []string{"internal", "core", "domain.go", "hoge.go"}
	if len(act) != len(wantNames) {
		t.Fatalf("want: %d, act: %d", len(wantNames), len(act))
	}
	for i, name := range wantNames {
		if act[i].Name != name {
			t.Errorf("want: %s, act[%d].Name: %s", name, i, act[i].Name)
		}
	}
}

func Test_Flatten_NestedCollapsed(t *testing.T) {
	t.Parallel()

	var act []*Node
	dir := &Node{
		Name:     "root",
		Path:     "/path",
		IsDir:    true,
		Expanded: true,
		Children: []*Node{
			{
				Name:     "internal",
				Path:     "/path/internal",
				IsDir:    true,
				Expanded: false,
				Children: []*Node{
					{Name: "hoge.go", Path: "/path/internal/hoge.go"},
				},
			},
		},
	}

	dir.Flatten(&act)

	wantNames := []string{"root", "internal"}
	if len(act) != len(wantNames) {
		t.Fatalf("want: %d, act: %d", len(wantNames), len(act))
	}
	for i, name := range wantNames {
		if act[i].Name != name {
			t.Errorf("want: %s, act[%d].Name: %s", name, i, act[i].Name)
		}
	}
}
