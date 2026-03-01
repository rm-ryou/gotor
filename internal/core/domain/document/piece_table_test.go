package document

import (
	"reflect"
	"testing"
)

func Test_Lines(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		text string
		want []string
	}

	testCases := []testCase{
		{
			name: "empty text",
			text: "",
			want: []string{""},
		},
		{
			name: "single line",
			text: "hello",
			want: []string{"hello"},
		},
		{
			name: "multiple lines",
			text: "a\nb",
			want: []string{"a", "b"},
		},
		{
			name: "trailing newline is ignored",
			text: "a\nb\n",
			want: []string{"a", "b"},
		},
		{
			name: "blank line is preserved",
			text: "a\n\n",
			want: []string{"a", ""},
		},
		{
			name: "single newline becomes one blank line",
			text: "\n",
			want: []string{""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pt := newPieceTableFromText(tc.text)
			got := pt.Lines()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Lines() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
