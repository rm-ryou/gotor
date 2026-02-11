package explorer

import "testing"

func Test_Flatten(t *testing.T) {
	type testCase struct {
		name      string
		node      *Node
		wantNames []string
	}

	testCases := []testCase{
		{
			name: "File node returns itself",
			node: &Node{
				Name:  "main.go",
				Path:  "/path/main.go",
				IsDir: false,
				Depth: 0,
			},
			wantNames: []string{"main.go"},
		},
		{
			name: "Collapsed directory returns only itself",
			node: &Node{
				Name:     "internal",
				Path:     "/path/internal",
				IsDir:    true,
				Depth:    0,
				Expanded: false,
				Children: []*Node{
					{Name: "hoge.go", Path: "/path/internal/hoge.go", Depth: 1},
				},
			},
			wantNames: []string{"internal"},
		},
		{
			name: "Expanded directory returns itself and child nodes",
			node: &Node{
				Name:     "internal",
				Path:     "path/internal",
				IsDir:    true,
				Depth:    0,
				Expanded: true,
				Children: []*Node{
					{
						Name:     "core",
						Path:     "/path/internal/core",
						IsDir:    true,
						Depth:    1,
						Expanded: true,
						Children: []*Node{
							{Name: "domain.go", Path: "/path/internal/core/domain.go", Depth: 2},
						},
					},
					{Name: "hoge.go", Path: "/path/internal/hoge.go", Depth: 1},
				},
			},
			wantNames: []string{"internal", "core", "domain.go", "hoge.go"},
		},
		{
			name: "If collapsed midway, nodes below are excluded",
			node: &Node{
				Name:     "root",
				Path:     "/path",
				IsDir:    true,
				Depth:    0,
				Expanded: true,
				Children: []*Node{
					{
						Name:     "internal",
						Path:     "/path/internal",
						IsDir:    true,
						Depth:    1,
						Expanded: false,
						Children: []*Node{
							{Name: "hoge.go", Path: "/path/internal/hoge.go", Depth: 2},
						},
					},
				},
			},
			wantNames: []string{"root", "internal"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var act []*Node

			tc.node.Flatten(&act)
			if len(act) != len(tc.wantNames) {
				t.Fatalf("want: %d, act: %d", len(tc.wantNames), len(act))
			}
			for i, name := range tc.wantNames {
				if act[i].Name != name {
					t.Errorf("want: %s, act[%d].Name: %s", name, i, act[i].Name)
				}
			}
		})
	}
}
