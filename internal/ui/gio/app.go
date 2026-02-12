package gio

import (
	"gioui.org/app"
	"gioui.org/unit"
)

type App struct{}

func New() *App {
	return &App{}
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
	w.Option(
		app.Title("gotor"),
		app.Size(unit.Dp(1200), unit.Dp(800)),
	)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		}
	}
}
