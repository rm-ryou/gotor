package document

import (
	"strings"
)

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

	if strings.HasSuffix(text, "\n") {
		text = strings.TrimSuffix(text, "\n")
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

func (pt *PieceTable) Insert(offset int, text string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if text == "" {
		return
	}

	totalLen := pt.textLen()
	if offset < 0 {
		offset = 0
	}
	if offset > totalLen {
		offset = totalLen
	}

	newPiece := Piece{
		bufferType: Add,
		start:      len(pt.addBuffer),
		length:     len(text),
	}
	pt.addBuffer += text

	if len(pt.pieces) == 0 {
		pt.pieces = []Piece{newPiece}
		return
	}

	newPieces := make([]Piece, 0, len(pt.pieces)+2)
	pos := 0
	inserted := false

	for _, piece := range pt.pieces {
		pieceEnd := pos + piece.length
		if !inserted && offset <= pieceEnd {
			rel := offset - pos
			switch {
			case rel <= 0:
				newPieces = append(newPieces, newPiece, piece)
			case rel >= piece.length:
				newPieces = append(newPieces, piece, newPiece)
			default:
				newPieces = appendPiece(newPieces, Piece{
					bufferType: piece.bufferType,
					start:      piece.start,
					length:     rel,
				})
				newPieces = append(newPieces, newPiece)
				newPieces = appendPiece(newPieces, Piece{
					bufferType: piece.bufferType,
					start:      piece.start + rel,
					length:     piece.length - rel,
				})
			}
			inserted = true
		} else {
			newPieces = append(newPieces, piece)
		}
		pos = pieceEnd
	}

	if !inserted {
		newPieces = append(newPieces, newPiece)
	}

	pt.pieces = newPieces
}

func (pt *PieceTable) Delete(start, end int) {
	if start < 0 {
		start = 0
	}

	totalLen := pt.textLen()
	if end > totalLen {
		end = totalLen
	}
	if start >= end {
		return
	}

	newPieces := make([]Piece, 0, len(pt.pieces))
	pos := 0

	for _, piece := range pt.pieces {
		pieceStart := pos
		pieceEnd := pos + piece.length

		if end <= pieceStart || start >= pieceEnd {
			newPieces = append(newPieces, piece)
			pos = pieceEnd
			continue
		}

		leftLen := start - pieceStart
		if leftLen > 0 {
			newPieces = appendPiece(newPieces, Piece{
				bufferType: piece.bufferType,
				start:      piece.start,
				length:     leftLen,
			})
		}

		rightLen := pieceEnd - end
		if rightLen > 0 {
			newPieces = appendPiece(newPieces, Piece{
				bufferType: piece.bufferType,
				start:      piece.start + piece.length - rightLen,
				length:     rightLen,
			})
		}

		pos = pieceEnd
	}

	pt.pieces = newPieces
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

func (pt *PieceTable) textLen() int {
	total := 0
	for _, piece := range pt.pieces {
		total += piece.length
	}
	return total
}

func appendPiece(pieces []Piece, piece Piece) []Piece {
	if piece.length <= 0 {
		return pieces
	}
	return append(pieces, piece)
}
