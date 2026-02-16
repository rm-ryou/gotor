package system

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const DefaultTextSize = 16

type Theme struct {
	*material.Theme
}

func NewTheme() (*Theme, error) {
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

	t.Theme.TextSize = unit.Sp(DefaultTextSize)
	t.Theme.Palette.Bg = color.NRGBA{R: 37, G: 37, B: 38, A: 255}
	t.Theme.Palette.Fg = color.NRGBA{R: 204, G: 204, B: 204, A: 255}

	return t, nil
}
