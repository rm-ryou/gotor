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

func TestPieceTableInsert(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		seed   string
		offset int
		text   string
		want   string
	}

	testCases := []testCase{
		{name: "insert into empty", seed: "", offset: 0, text: "abc", want: "abc"},
		{name: "insert at beginning", seed: "world", offset: 0, text: "hello ", want: "hello world"},
		{name: "insert in middle", seed: "helo", offset: 2, text: "l", want: "hello"},
		{name: "insert at end", seed: "hello", offset: 5, text: "!", want: "hello!"},
		{name: "normalize line endings", seed: "a", offset: 1, text: "\r\nb", want: "a\nb"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pt := newPieceTableFromText(tc.seed)
			pt.Insert(tc.offset, tc.text)

			if act := pt.Text(); act != tc.want {
				t.Fatalf("Text() = %q, want %q", act, tc.want)
			}
		})
	}
}

func TestPieceTableDelete(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		seed  string
		start int
		end   int
		want  string
	}

	testCases := []testCase{
		{name: "delete from start", seed: "hello", start: 0, end: 2, want: "llo"},
		{name: "delete from middle", seed: "hello", start: 1, end: 4, want: "ho"},
		{name: "delete from end", seed: "hello", start: 3, end: 5, want: "hel"},
		{name: "delete whole text", seed: "hello", start: 0, end: 5, want: ""},
		{name: "clamp range", seed: "hello", start: -10, end: 99, want: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pt := newPieceTableFromText(tc.seed)
			pt.Delete(tc.start, tc.end)

			if act := pt.Text(); act != tc.want {
				t.Fatalf("Text() = %q, want %q", act, tc.want)
			}
		})
	}
}
