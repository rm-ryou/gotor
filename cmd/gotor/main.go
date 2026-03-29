package main

import (
	"fmt"
	"os"

	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/platform/fs"
	"github.com/rm-ryou/gotor/internal/ui/gio/app"
	"github.com/rm-ryou/gotor/internal/ui/gio/view"
)

func main() {
	// TODO: refactor
	fsIO := fs.New(true)
	explorerUC, err := usecase.NewExplorer(fsIO, fsIO, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	view, err := view.New(explorerUC)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a := app.New(view)

	if err := a.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
