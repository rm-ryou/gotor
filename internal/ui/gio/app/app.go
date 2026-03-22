package app

import (
	"image"
	"image/color"
	"time"

	"gioui.org/app"
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

type App struct {
	theme  *system.Theme
	window *app.Window

	explorerView *features.ExplorerView
	editorView   *features.EditorView

	errorPopup errorPopup
}

type errorPopup struct {
	message   string
	expiresAt time.Time
}

func New(explorerUC *usecase.Explorer) (*App, error) {
	w := new(app.Window)
	w.Option(app.Title("gotor"))
	th, err := system.NewTheme()
	if err != nil {
		return nil, err
	}

	explorerView := features.NewExplorerView(th, explorerUC)

	fileIO := fs.NewFileIO()
	editorUC := usecase.NewEditor(fileIO)

	explorerUC.OnFileSelected = func(path string) error {
		return editorUC.OpenFile(path)
	}
	editorView := features.NewEditorView(th, editorUC)

	a := &App{
		theme:  th,
		window: w,

		explorerView: explorerView,
		editorView:   editorView,
	}
	a.explorerView.OnError = a.showError
	a.editorView.OnError = a.showError

	return a, nil
}

func (a *App) Run() error {
	errCh := make(chan error)

	go func() {
		errCh <- a.loop()
	}()
	app.Main()

	return <-errCh
}

func (a *App) loop() error {
	var ops op.Ops
	for {
		switch e := a.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.setEventHandler(gtx)
			a.layout(gtx, a.theme)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) layout(gtx layout.Context, th *system.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, th.Palette.Bg)

	dims := layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(explorerPaneWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width

			return a.explorerView.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := gtx.Dp(unit.Dp(paneDividerWidth))
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, th.Palette.Fg)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.editorView.Layout(gtx)
		}),
	)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return a.layoutErrorPopup(gtx)
		}),
	)
}

func (a *App) setEventHandler(gtx layout.Context) {
	a.explorerView.HandleNodeClicks(gtx)
	a.editorView.HandleKeyInput(gtx)
}

func (a *App) showError(err error) {
	a.errorPopup = errorPopup{
		message:   usecase.MessageFor(err),
		expiresAt: time.Now().Add(4 * time.Second),
	}
	a.window.Invalidate()
}

func (a *App) layoutErrorPopup(gtx layout.Context) layout.Dimensions {
	if a.errorPopup.message == "" {
		return layout.Dimensions{}
	}
	if gtx.Now.After(a.errorPopup.expiresAt) {
		a.errorPopup = errorPopup{}
		return layout.Dimensions{}
	}

	gtx.Execute(op.InvalidateCmd{At: a.errorPopup.expiresAt})

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
						lbl := material.Body2(a.theme.Theme, a.errorPopup.message)
						lbl.Color = color.NRGBA{R: 255, G: 244, B: 246, A: 255}
						lbl.WrapPolicy = text.WrapWords
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	})
}
