package main

import (
	"fmt"
	"os"

	"github.com/rm-ryou/gotor/internal/ui/gio"
)

func main() {
	app := gio.New()
	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
