package system

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	assetfonts "github.com/rm-ryou/gotor/internal/ui/assets/fonts"
)

const DefaultTextSize = 12

type Theme struct {
	*material.Theme
}

func NewTheme() (*Theme, error) {
	th := material.NewTheme()
	fontFaces, err := Prepare()
	if err != nil {
		return nil, err
	}

	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontFaces))
	th.Face = font.Typeface(assetfonts.DefaultTypeface)

	t := &Theme{
		Theme: th,
	}

	t.Theme.TextSize = unit.Sp(DefaultTextSize)
	t.Theme.Palette.Bg = color.NRGBA{R: 37, G: 37, B: 38, A: 255}
	t.Theme.Palette.Fg = color.NRGBA{R: 204, G: 204, B: 204, A: 255}

	return t, nil
}
