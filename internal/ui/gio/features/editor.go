package features

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
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
	mode  EditorDisplayMode

	keyTag  struct{}
	focused bool

	OnError func(error)
}

type EditorDisplayMode uint8

const (
	EditorDisplayModeWrap EditorDisplayMode = iota
	EditorDisplayModeHorizontalScroll
)

const (
	tabWidth          = 4
	editorInsetTop    = unit.Dp(8)
	editorInsetBottom = unit.Dp(8)
	editorInsetLeft   = unit.Dp(8)
	lineNumberGap     = unit.Dp(10)
	cursorLineHeight  = unit.Dp(22)
	cursorStrokeWidth = unit.Dp(1)
)

func NewEditorView(th *system.Theme, uc *usecase.Editor) *EditorView {
	return &EditorView{
		theme: th,
		uc:    uc,
		lines: []string{},
		mode:  EditorDisplayModeHorizontalScroll,
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
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return ev.layoutCursor(gtx)
		}),
	)
}

func (ev *EditorView) HandleKeyInput(gtx layout.Context) {
	event.Op(gtx.Ops, &ev.keyTag)
	key.InputHintOp{Tag: &ev.keyTag, Hint: key.HintText}.Add(gtx.Ops)
	if !ev.focused {
		gtx.Execute(key.FocusCmd{Tag: &ev.keyTag})
	}

	for {
		ke, ok := gtx.Event(
			key.FocusFilter{Target: &ev.keyTag},
			key.Filter{Focus: &ev.keyTag, Name: key.NameLeftArrow},
			key.Filter{Focus: &ev.keyTag, Name: key.NameRightArrow},
			key.Filter{Focus: &ev.keyTag, Name: key.NameUpArrow},
			key.Filter{Focus: &ev.keyTag, Name: key.NameDownArrow},
			key.Filter{Focus: &ev.keyTag, Name: key.NameDeleteBackward},
			key.Filter{Focus: &ev.keyTag, Name: key.NameEnter},
			key.Filter{Focus: &ev.keyTag, Name: key.NameReturn},
			key.Filter{Focus: &ev.keyTag, Name: "S", Required: key.ModShortcut},
		)
		if !ok {
			break
		}

		switch ke := ke.(type) {
		case key.FocusEvent:
			ev.focused = ke.Focus
		case key.Event:
			if ke.State != key.Press {
				continue
			}
			if ev.handleKeyEvent(ke) {
				ev.ensureCursorVisible()
			}
		case key.EditEvent:
			if ev.handleTextInput(ke.Text) {
				ev.ensureCursorVisible()
			}
		}
	}
}

func (ev *EditorView) layoutContent(gtx layout.Context, lineWidth int, lines []string, textColor color.NRGBA) layout.Dimensions {
	return layout.Inset{
		Top: editorInsetTop, Bottom: editorInsetBottom,
		Left: editorInsetLeft,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		itemCount := len(lines) + 1
		return material.List(ev.theme.Theme, &ev.list).Layout(
			gtx, itemCount,
			func(gtx layout.Context, i int) layout.Dimensions {
				if i == len(lines) {
					h := gtx.Dp(cursorLineHeight)
					return layout.Dimensions{Size: image.Pt(0, h)}
				}
				return ev.layoutLine(gtx, lineWidth, i+1, lines[i], textColor)
			},
		)
	})
}

func (ev *EditorView) layoutLine(gtx layout.Context, lineWidth, lineNum int, lineText string, textColor color.NRGBA) layout.Dimensions {
	lineNumColor := color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	displayText := expandTabs(lineText, tabWidth)
	rowHeight := gtx.Dp(cursorLineHeight)

	gtx.Constraints.Min.Y = rowHeight
	gtx.Constraints.Max.Y = rowHeight

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			minWidth := lineWidth
			gtx.Constraints.Min.X = minWidth

			return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(ev.theme.Theme, strconv.Itoa(lineNum))
				lbl.Color = lineNumColor
				lbl.Alignment = text.End
				lbl.MaxLines = 1
				return lbl.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(ev.theme.Theme, displayText)
			lbl.Color = textColor
			lbl.MaxLines = 1
			return lbl.Layout(gtx)
		}),
	)
}

func (ev *EditorView) layoutCursor(gtx layout.Context) layout.Dimensions {
	cursor := ev.uc.Cursor()
	lines := ev.uc.Document().Lines()
	if len(lines) == 0 || cursor.Row >= len(lines) {
		return layout.Dimensions{}
	}

	cursorColor := color.NRGBA{R: 120, G: 200, B: 255, A: 180}
	lineNumWidth := gtx.Dp(unit.Dp(10)) * len(strconv.Itoa(len(lines)))
	lineHeight := gtx.Dp(cursorLineHeight)
	leftPadding := gtx.Dp(editorInsetLeft)
	lineNumberPadding := gtx.Dp(lineNumberGap)
	cursorWidth := gtx.Dp(cursorStrokeWidth)
	cursorHeight := gtx.Sp(ev.theme.TextSize)

	rowOffset := cursor.Row - ev.list.Position.First
	x := leftPadding + lineNumWidth + lineNumberPadding +
		ev.measureTextWidth(gtx, displayPrefixForCursor(lines[cursor.Row], cursor.Col, tabWidth))
	y := gtx.Dp(editorInsetTop) - ev.list.Position.Offset + (rowOffset * lineHeight) + (lineHeight-cursorHeight)/2 - 1
	contentRect := image.Rectangle{
		Min: image.Point{
			X: leftPadding + lineNumWidth + lineNumberPadding,
			Y: gtx.Dp(editorInsetTop),
		},
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y - gtx.Dp(editorInsetBottom),
		},
	}

	cursorRect := image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + cursorWidth, Y: y + cursorHeight},
	}

	defer clip.Rect(contentRect).Push(gtx.Ops).Pop()
	defer clip.Rect(cursorRect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, cursorColor)

	return layout.Dimensions{
		Size: image.Point{X: cursorWidth, Y: cursorHeight},
	}
}

func displayColumnForCursor(s string, limit, width int) int {
	if limit <= 0 {
		return 0
	}

	column := 0
	index := 0

	for _, r := range s {
		if index >= limit {
			break
		}

		if r == '\t' {
			column += width - (column % width)
		} else {
			column++
		}

		index++
	}

	return column
}

func displayPrefixForCursor(s string, limit, width int) string {
	if limit <= 0 {
		return ""
	}
	return expandTabs(string([]rune(s)[:limit]), width)
}

func (ev *EditorView) measureTextWidth(gtx layout.Context, value string) int {
	var ops op.Ops
	lbl := material.Body2(ev.theme.Theme, value)
	lbl.MaxLines = 1
	return lbl.Layout(layout.Context{Constraints: layout.Constraints{Max: image.Pt(1_000_000, 1_000_000)}, Metric: gtx.Metric, Now: gtx.Now, Locale: gtx.Locale, Values: gtx.Values, Ops: &ops}).Size.X
}

func (ev *EditorView) handleKeyEvent(evt key.Event) bool {
	if evt.Name == "S" && evt.Modifiers.Contain(key.ModShortcut) {
		if err := ev.uc.Save(); err != nil {
			ev.reportError(err)
		}
		return false
	}

	switch evt.Name {
	case key.NameUpArrow:
		ev.uc.MoveCursorUp()
	case key.NameDownArrow:
		ev.uc.MoveCursorDown()
	case key.NameLeftArrow:
		ev.uc.MoveCursorLeft()
	case key.NameRightArrow:
		ev.uc.MoveCursorRight()
	case key.NameDeleteBackward:
		return ev.uc.DeleteBackward()
	case key.NameEnter, key.NameReturn:
		ev.uc.InsertText("\n")
	default:
		return false
	}

	return true
}

func (ev *EditorView) reportError(err error) {
	if err == nil || ev.OnError == nil {
		return
	}
	ev.OnError(err)
}

func (ev *EditorView) handleTextInput(text string) bool {
	if text == "" {
		return false
	}

	ev.uc.InsertText(text)
	return true
}

func (ev *EditorView) ensureCursorVisible() {
	cursorRow := ev.uc.Cursor().Row

	if cursorRow < ev.list.Position.First {
		ev.list.ScrollTo(cursorRow)
		return
	}

	if ev.list.Position.Count <= 0 {
		return
	}

	lastVisibleRow := ev.list.Position.First + ev.list.Position.Count - 1
	if cursorRow > lastVisibleRow {
		ev.list.ScrollTo(cursorRow - ev.list.Position.Count + 1)
	}
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
