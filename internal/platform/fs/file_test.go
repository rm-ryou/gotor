package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Read(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name    string
		content string
		setup   func(t *testing.T, content string) string
		wantErr bool
	}

	testCases := []testCase{
		{
			name: "when path is exists. return file contents",
			setup: func(t *testing.T, content string) string {
				t.Helper()

				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "test.txt")
				if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return testFile
			},
			wantErr: false,
		},
		{
			name: "when path is not exists, return err",
			setup: func(t *testing.T, content string) string {
				t.Helper()
				return "/nonexistent/file/path.txt"
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := tc.setup(t, tc.content)

			fio := NewFileIO()
			got, err := fio.Read(path)

			if (err != nil) != tc.wantErr {
				t.Errorf("want error: %t, act: %v", tc.wantErr, err)
			}

			if got != tc.content {
				t.Errorf("want content: %s, act: %s", got, tc.content)
			}
		})
	}
}
