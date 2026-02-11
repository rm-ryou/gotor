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
