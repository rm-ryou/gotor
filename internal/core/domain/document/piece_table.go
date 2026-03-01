package document

import "strings"

type BufferType int

const (
	Org BufferType = iota
	Add
)

type Piece struct {
	bufferType BufferType
	start      int
	length     int
}

type PieceTable struct {
	orgBuffer string
	addBuffer string
	pieces    []Piece
}

func newPieceTable() *PieceTable {
	return &PieceTable{
		orgBuffer: "",
		addBuffer: "",
		pieces:    []Piece{},
	}
}

func newPieceTableFromText(text string) *PieceTable {
	text = strings.ReplaceAll(text, "\r\n", "\n")

	if text == "" {
		return &PieceTable{
			orgBuffer: "",
			addBuffer: "",
			pieces:    []Piece{},
		}
	}

	return &PieceTable{
		orgBuffer: text,
		addBuffer: "",
		pieces: []Piece{
			{
				bufferType: Org,
				start:      0,
				length:     len(text),
			},
		},
	}
}

func (pt *PieceTable) Lines() []string {
	text := pt.Text()
	if text == "" {
		return []string{""}
	}
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func (pt *PieceTable) LineCount() int {
	return len(pt.Lines())
}

func (pt *PieceTable) Line(row int) string {
	lines := pt.Lines()
	if row < 0 || row >= len(lines) {
		return ""
	}
	return lines[row]
}

func (pt *PieceTable) Text() string {
	var sb strings.Builder
	for _, piece := range pt.pieces {
		sb.WriteString(pt.getPieceText(piece))
	}
	return sb.String()
}

func (pt *PieceTable) getPieceText(piece Piece) string {
	switch piece.bufferType {
	case Org:
		return pt.orgBuffer[piece.start : piece.start+piece.length]
	case Add:
		return pt.addBuffer[piece.start : piece.start+piece.length]
	default:
		return ""
	}
}
