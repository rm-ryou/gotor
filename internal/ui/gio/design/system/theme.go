package system

import (
	"gioui.org/font"
	"gioui.org/widget/material"
	"github.com/rm-ryou/gotor/internal/ui/gio/config"
)

type Theme struct {
	*material.Theme
}

func NewTheme(cfg config.Theme) (*Theme, error) {
	th := material.NewTheme()
	fontFaces, err := Prepare()
	if err != nil {
		return nil, err
	}

	th.Shaper = NewShaper(fontFaces)
	th.Face = font.Typeface(DefaultTypefaceWithFallback())

	t := &Theme{
		Theme: th,
	}

	t.Theme.TextSize = cfg.TextSize
	t.Theme.Palette.Bg = cfg.Bg
	t.Theme.Palette.Fg = cfg.Fg

	return t, nil
}
