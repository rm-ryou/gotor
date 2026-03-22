package view

import (
	"image"
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
	"github.com/rm-ryou/gotor/internal/ui/gio/config"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
	"github.com/rm-ryou/gotor/internal/ui/gio/features"
)

type View struct {
	theme *system.Theme
	cfg   config.UI

	explorer *features.ExplorerView
	editor   *features.EditorView

	errorPopup errorPopup
}

type errorPopup struct {
	message   string
	expiresAt time.Time
}

func New(explorerUC *usecase.Explorer) (*View, error) {
	return NewWithConfig(explorerUC, config.Default())
}

func NewWithConfig(explorerUC *usecase.Explorer, cfg config.UI) (*View, error) {
	th, err := system.NewTheme(cfg.Theme)
	if err != nil {
		return nil, err
	}

	fileIO := fs.NewFileIO()
	editorUC := usecase.NewEditor(fileIO)

	view := &View{
		theme:    th,
		cfg:      cfg,
		explorer: features.NewExplorerView(th, explorerUC, cfg.Explorer),
		editor:   features.NewEditorView(th, editorUC, cfg.Editor),
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
			width := gtx.Dp(unit.Dp(v.cfg.View.ExplorerPaneWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width

			return v.explorer.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(v.cfg.View.PaneDividerWidth))
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
		expiresAt: time.Now().Add(v.cfg.View.ErrorPopup.Duration),
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

	inset := layout.UniformInset(v.cfg.View.ErrorPopup.Inset)
	return layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			maxWidth := gtx.Dp(v.cfg.View.ErrorPopup.MaxWidth)
			if gtx.Constraints.Max.X > maxWidth {
				gtx.Constraints.Max.X = maxWidth
			}

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(v.cfg.View.ErrorPopup.Radius)).Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, v.cfg.View.ErrorPopup.Bg)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    v.cfg.View.ErrorPopup.PaddingY,
						Right:  v.cfg.View.ErrorPopup.PaddingX,
						Bottom: v.cfg.View.ErrorPopup.PaddingY,
						Left:   v.cfg.View.ErrorPopup.PaddingX,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(v.theme.Theme, v.errorPopup.message)
						lbl.Color = v.cfg.View.ErrorPopup.TextColor
						lbl.WrapPolicy = text.WrapWords
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	})
}
