package layout

import "gioui.org/unit"

const DefaultIndentPerDepth = 12

type Explorer struct {
	indentPerDepth int
	rowHeight      int
}

func NewExplorer(textSize, indentPerDepth, rowHeightDelta int) *Explorer {
	return &Explorer{
		indentPerDepth: indentPerDepth,
		rowHeight:      textSize + rowHeightDelta,
	}
}

func (e *Explorer) Indent(depth int) unit.Dp {
	return unit.Dp(float32(depth) * float32(e.indentPerDepth))
}

func (e *Explorer) RowHeight() unit.Dp {
	return unit.Dp(e.rowHeight)
}
