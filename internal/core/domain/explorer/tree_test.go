package explorer

import (
	"testing"
)

func TestTree_VisibleNodes_InitialState(t *testing.T) {
	t.Parallel()

	rootNode := &Node{
		Name:     "root",
		Path:     "/path",
		IsDir:    true,
		Depth:    0,
		Expanded: true,
		Children: []*Node{
			{Name: "main.go", Path: "/path/main.go", Depth: 1},
			{
				Name:     "internal",
				Path:     "/path/internal",
				IsDir:    true,
				Depth:    1,
				Expanded: true,
				Children: []*Node{
					{Name: "core.go", Path: "/path/internal/core.go", Depth: 2},
				},
			},
		},
	}
	tree := New(rootNode)

	act := tree.VisibleNodes()

	wantNames := []string{"main.go", "internal", "core.go"}
	if len(act) != len(wantNames) {
		t.Fatalf("want: %d, act: %d", len(wantNames), len(act))
	}
	for i, name := range wantNames {
		if act[i].Name != name {
			t.Errorf("want: %s, act[%d].Name: %s", name, i, act[i].Name)
		}
	}
}

func TestTree_VisibleNodes_CachedOnSecondCall(t *testing.T) {
	t.Parallel()

	rootNode := &Node{
		Name:     "root",
		Path:     "/path",
		IsDir:    true,
		Depth:    0,
		Expanded: true,
		Children: []*Node{
			{Name: "main.go", Path: "/path/main.go", Depth: 1},
		},
	}
	tree := New(rootNode)

	first := tree.VisibleNodes()
	second := tree.VisibleNodes()

	if &first[0] != &second[0] {
		t.Error("expected same slice on second call (cache hit)")
	}
}

func Test_Select(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name             string
		target           *Node
		wantSelectedPath string
	}

	testCases := []testCase{
		{
			name:             "set file name, when file selected",
			target:           &Node{Name: "main.go", Path: "path/main.go", Depth: 1},
			wantSelectedPath: "path/main.go",
		},
		{
			name:             "nothing to do, when directory selected",
			target:           &Node{Name: "internal", Path: "path/internal", IsDir: true, Depth: 1},
			wantSelectedPath: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rootNode := &Node{
				Name:     "root",
				Path:     "/path",
				IsDir:    true,
				Depth:    0,
				Expanded: true,
				Children: []*Node{tc.target},
			}

			tree := New(rootNode)
			tree.Select(tc.target)

			if tree.SelectedPath() != tc.wantSelectedPath {
				t.Errorf("want: %s, act: %s", tc.wantSelectedPath, tree.SelectedPath())
			}
		})
	}
}

func Test_FindNode(t *testing.T) {
	t.Parallel()

	file := &Node{
		Name:     "hoge.go",
		Path:     "/path/internal/hoge.go",
		IsDir:    false,
		Depth:    2,
		Expanded: false,
	}

	dir := &Node{
		Name:     "internal",
		Path:     "/path/internal",
		IsDir:    true,
		Depth:    1,
		Expanded: true,
		Children: []*Node{file},
	}

	rootNode := &Node{
		Name:     "root",
		Path:     "/path",
		IsDir:    true,
		Depth:    0,
		Expanded: true,
		Children: []*Node{dir},
	}

	tree := New(rootNode)

	type testCase struct {
		path string
		want *Node
	}

	testCases := []testCase{
		{path: "/path", want: rootNode},
		{path: "/path/internal", want: dir},
		{path: "/path/internal/hoge.go", want: file},
		{path: "/path/notexist.go", want: nil},
		{path: "/other/path/hoge.go", want: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			act := tree.FindNode(tc.path)

			if act != tc.want {
				t.Errorf("want: %v, act: %v", tc.want, act)
			}
		})
	}
}
