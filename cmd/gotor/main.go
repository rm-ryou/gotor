package main

import (
	"fmt"
	"os"

	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/platform/fs"
	"github.com/rm-ryou/gotor/internal/ui/gio/app"
)

func main() {
	// TODO: refactor
	fs := fs.New(true)
	explorerUC, err := usecase.NewExplorer(fs, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a, err := app.New(explorerUC)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
