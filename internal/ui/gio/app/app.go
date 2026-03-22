package app

import (
	"gioui.org/app"
	"gioui.org/op"
	"github.com/rm-ryou/gotor/internal/ui/gio/view"
)

type App struct {
	window *app.Window
	view   *view.View
}

func New(v *view.View) *App {
	w := new(app.Window)
	w.Option(app.Title("gotor"))

	return &App{
		window: w,
		view:   v,
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
	var ops op.Ops
	for {
		switch e := a.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.view.HandleEvents(gtx)
			a.view.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
