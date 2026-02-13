package app

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
	"github.com/rm-ryou/gotor/internal/ui/gio/features"
)

type App struct {
	theme  *system.Theme
	window *app.Window

	explorerView *features.ExplorerView
}

func New(explorerUC *usecase.Explorer) *App {
	w := new(app.Window)
	th := system.NewTheme()

	explorerView := features.NewExplorerView(th, explorerUC)

	return &App{
		theme:  th,
		window: w,

		explorerView: explorerView,
	}
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
	w := new(app.Window)
	w.Option(app.Title("gotor"))

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.Layout(gtx, a.theme)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) Layout(gtx layout.Context, th *system.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, th.Palette.Bg)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
