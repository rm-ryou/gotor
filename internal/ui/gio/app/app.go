package app

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/rm-ryou/gotor/internal/core/usecase"
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
}

func New(explorerUC *usecase.Explorer) (*App, error) {
	w := new(app.Window)
	w.Option(app.Title("gotor"))
	th, err := system.NewTheme()
	if err != nil {
		return nil, err
	}

	explorerView := features.NewExplorerView(th, explorerUC)
	editorView := features.NewEditorView(th)

	return &App{
		theme:  th,
		window: w,

		explorerView: explorerView,
		editorView:   editorView,
	}, nil
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

	return dims
}

func (a *App) setEventHandler(gtx layout.Context) {
	a.explorerView.HandleNodeClicks(gtx)
}
