package view

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/platform/fs"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
	"github.com/rm-ryou/gotor/internal/ui/gio/features"
)

const (
	explorerPaneWidth = 200
	paneDividerWidth  = 1
)

type View struct {
	theme *system.Theme

	explorer *features.ExplorerView
	editor   *features.EditorView

	errorPopup errorPopup
}

type errorPopup struct {
	message   string
	expiresAt time.Time
}

func New(explorerUC *usecase.Explorer) (*View, error) {
	th, err := system.NewTheme()
	if err != nil {
		return nil, err
	}

	fileIO := fs.NewFileIO()
	editorUC := usecase.NewEditor(fileIO)

	view := &View{
		theme:    th,
		explorer: features.NewExplorerView(th, explorerUC),
		editor:   features.NewEditorView(th, editorUC),
	}

	explorerUC.OnFileSelected = func(path string) error {
		return editorUC.OpenFile(path)
	}
	view.explorer.OnError = view.showError
	view.editor.OnError = view.showError

	return view, nil
}

func (v *View) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, v.theme.Palette.Bg)

	content := layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(explorerPaneWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width

			return v.explorer.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(paneDividerWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, v.theme.Palette.Fg)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return v.editor.Layout(gtx)
		}),
	)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return content
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return v.layoutErrorPopup(gtx)
		}),
	)
}

func (v *View) HandleEvents(gtx layout.Context) {
	v.explorer.HandleNodeClicks(gtx)
	v.editor.HandleKeyInput(gtx)
}

func (v *View) showError(err error) {
	v.errorPopup = errorPopup{
		message:   usecase.MessageFor(err),
		expiresAt: time.Now().Add(4 * time.Second),
	}
}

func (v *View) layoutErrorPopup(gtx layout.Context) layout.Dimensions {
	if v.errorPopup.message == "" {
		return layout.Dimensions{}
	}
	if gtx.Now.After(v.errorPopup.expiresAt) {
		v.errorPopup = errorPopup{}
		return layout.Dimensions{}
	}

	gtx.Execute(op.InvalidateCmd{At: v.errorPopup.expiresAt})

	inset := layout.UniformInset(unit.Dp(16))
	return layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			maxWidth := gtx.Dp(unit.Dp(320))
			if gtx.Constraints.Max.X > maxWidth {
				gtx.Constraints.Max.X = maxWidth
			}

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(unit.Dp(10))).Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, color.NRGBA{R: 177, G: 55, B: 72, A: 235})
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    unit.Dp(12),
						Right:  unit.Dp(16),
						Bottom: unit.Dp(12),
						Left:   unit.Dp(16),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(v.theme.Theme, v.errorPopup.message)
						lbl.Color = color.NRGBA{R: 255, G: 244, B: 246, A: 255}
						lbl.WrapPolicy = text.WrapWords
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	})
}
