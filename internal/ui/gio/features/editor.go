package features

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

const editorLineNumberWidth = 48

var editorPlaceholderLines = []string{
	"hello, World!",
}

type EditorView struct {
	theme *system.Theme
}

func NewEditorView(th *system.Theme) *EditorView {
	return &EditorView{
		theme: th,
	}
}

func (ev *EditorView) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, ev.theme.Palette.Bg)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, editorLineWidgets(ev)...)
			})
		}),
	)
}

func (ev *EditorView) layoutLine(gtx layout.Context, lineNumber int, content string) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(editorLineNumberWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width

			lbl := material.Body2(ev.theme.Theme, fmt.Sprintf("%d", lineNumber))
			lbl.Color = ev.theme.Palette.Fg
			return lbl.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(ev.theme.Theme, content)
			lbl.Color = ev.theme.Palette.Fg
			return lbl.Layout(gtx)
		}),
	)
}

func editorLineWidgets(ev *EditorView) []layout.FlexChild {
	children := make([]layout.FlexChild, 0, len(editorPlaceholderLines))
	for i, line := range editorPlaceholderLines {
		lineNumber := i + 1
		content := line
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ev.layoutLine(gtx, lineNumber, content)
		}))
	}

	return children
}
