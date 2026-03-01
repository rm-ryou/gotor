package features

import (
	"image/color"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

type EditorView struct {
	theme *system.Theme
	uc    *usecase.Editor
	lines []string
	list  widget.List
}

const tabWidth = 4

func NewEditorView(th *system.Theme, uc *usecase.Editor) *EditorView {
	return &EditorView{
		theme: th,
		uc:    uc,
		lines: []string{},
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
}

func (ev *EditorView) Layout(gtx layout.Context) layout.Dimensions {
	textColor := color.NRGBA{R: 212, G: 212, B: 212, A: 255}
	gtx.Constraints.Min = gtx.Constraints.Max

	lines := ev.uc.Document().Lines()

	numLines := len(lines)
	lineWidth := gtx.Dp(unit.Dp(10)) * len(strconv.Itoa(numLines))

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, ev.theme.Palette.Bg)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return ev.layoutContent(gtx, lineWidth, lines, textColor)
		}),
	)
}

func (ev *EditorView) layoutContent(gtx layout.Context, lineWidth int, lines []string, textColor color.NRGBA) layout.Dimensions {
	return layout.Inset{
		Top: unit.Dp(8), Bottom: unit.Dp(8),
		Left: unit.Dp(8), Right: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(ev.theme.Theme, &ev.list).Layout(
			gtx, len(lines),
			func(gtx layout.Context, i int) layout.Dimensions {
				return ev.layoutLine(gtx, lineWidth, i+1, lines[i], textColor)
			},
		)
	})
}

func (ev *EditorView) layoutLine(gtx layout.Context, lineWidth, lineNum int, lineText string, textColor color.NRGBA) layout.Dimensions {
	lineNumColor := color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	displayText := expandTabs(lineText, tabWidth)

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Baseline,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			minWidth := lineWidth
			gtx.Constraints.Min.X = minWidth

			return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(ev.theme.Theme, strconv.Itoa(lineNum))
				lbl.Color = lineNumColor
				lbl.Alignment = text.End
				return lbl.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(ev.theme.Theme, displayText)
			lbl.Color = textColor
			return lbl.Layout(gtx)
		}),
	)
}

func expandTabs(s string, width int) string {
	if width <= 0 || !strings.ContainsRune(s, '\t') {
		return s
	}

	var b strings.Builder
	column := 0

	for _, r := range s {
		if r == '\t' {
			spaces := width - (column % width)
			for range spaces {
				b.WriteByte(' ')
			}
			column += spaces
			continue
		}

		b.WriteRune(r)
		column++
	}

	return b.String()
}
