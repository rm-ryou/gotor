package explorer

import (
	"errors"
	"testing"

	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
	"github.com/rm-ryou/gotor/internal/core/domain/explorer/mocks"
	"go.uber.org/mock/gomock"
)

func Test_New(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		rootPath      string
		setupMock     func(*mocks.MockFSReader)
		wantErr       bool
		checkRootPath bool
	}

	testCases := []testCase{
		{
			name:     "when rootPath is empty Initialize Getwd()",
			rootPath: "",
			setupMock: func(m *mocks.MockFSReader) {
				m.EXPECT().ReadDir(gomock.Any(), 1).Return(nil, nil)
			},
			wantErr:       false,
			checkRootPath: false,
		},
		{
			name:     "when err in ReadDir return err",
			rootPath: t.TempDir(),
			setupMock: func(m *mocks.MockFSReader) {
				m.EXPECT().ReadDir(gomock.Any(), 1).Return(nil, errors.New("permission denied"))
			},
			wantErr:       true,
			checkRootPath: false,
		},
		{
			name:          "when rootPath is not exists return err",
			rootPath:      "/nonexistent/path/",
			setupMock:     func(m *mocks.MockFSReader) {},
			wantErr:       true,
			checkRootPath: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockFS := mocks.NewMockFSReader(ctrl)
			tc.setupMock(mockFS)

			uc, err := New(mockFS, tc.rootPath)
			if (err != nil) != tc.wantErr {
				t.Errorf("want error: %t, act: %v", tc.wantErr, err)
			}
			if tc.wantErr {
				return
			}

			if tc.checkRootPath && uc.Tree().Root().Path != tc.rootPath {
				t.Errorf("want rootPath: %s, act: %s", tc.rootPath, uc.Tree().Root().Path)
			}
		})
	}
}

func Test_ToggleNode(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	subdir := root + "/internal"

	type testCase struct {
		name               string
		initialNodes       []*domain.Node
		subdirNodes        []*domain.Node
		toggleCalled       int
		readDirErr         error
		wantErr            bool
		wantExpanded       bool
		wantChildName      string
		wantReadDirsCalled int
	}

	testCases := []testCase{
		{
			name: "Call ReadDir on first expand",
			initialNodes: []*domain.Node{
				{Name: "internal", Path: subdir, IsDir: true, Depth: 1},
			},
			subdirNodes: []*domain.Node{
				{Name: "core.go", Path: subdir + "/core.go", IsDir: false, Depth: 2},
			},
			toggleCalled:       1,
			wantErr:            false,
			wantExpanded:       true,
			wantChildName:      "core.go",
			wantReadDirsCalled: 2,
		},
		{
			name: "Do not call ReadDir on second expand",
			initialNodes: []*domain.Node{
				{Name: "internal", Path: subdir, IsDir: true, Depth: 1},
			},
			subdirNodes: []*domain.Node{
				{Name: "core.go", Path: subdir + "/core.go", IsDir: false, Depth: 2},
			},
			toggleCalled:       3,
			wantErr:            false,
			wantExpanded:       true,
			wantChildName:      "core.go",
			wantReadDirsCalled: 2,
		},
		{
			name: "Hide children when collapsed",
			initialNodes: []*domain.Node{
				{Name: "internal", Path: subdir, IsDir: true, Depth: 1},
			},
			subdirNodes: []*domain.Node{
				{Name: "core.go", Path: subdir + "/core.go", IsDir: false, Depth: 2},
			},
			toggleCalled:       2,
			wantErr:            false,
			wantExpanded:       false,
			wantChildName:      "",
			wantReadDirsCalled: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockFS := mocks.NewMockFSReader(ctrl)
			mockFS.EXPECT().ReadDir(root, 1).Return(tc.initialNodes, nil)
			if tc.readDirErr != nil {
				mockFS.EXPECT().ReadDir(subdir, 2).Return(nil, tc.readDirErr)
			} else {
				mockFS.EXPECT().ReadDir(subdir, 2).Return(tc.subdirNodes, nil).MaxTimes(1)
			}

			uc, _ := New(mockFS, root)
			node := uc.Tree().VisibleNodes()[0]

			var err error
			for range tc.toggleCalled {
				err = uc.ToggleNode(node)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("want error: %t, act: %v", tc.wantErr, err)
			}
			if node.Expanded != tc.wantExpanded {
				t.Errorf("want expanded: %t, act: %t", tc.wantExpanded, node.Expanded)
			}

			if tc.wantChildName != "" {
				found := false
				for _, n := range uc.Tree().VisibleNodes() {
					if n.Name == tc.wantChildName {
						found = true
					}
				}

				if !found {
					t.Errorf("child node should include %s node", tc.wantChildName)
				}
			}
		})
	}
}
