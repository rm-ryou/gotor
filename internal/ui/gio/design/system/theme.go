package system

import (
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget/material"
)

const DefaultTextSize = 12

type Theme struct {
	*material.Theme
}

func NewTheme() *Theme {
	th := material.NewTheme()
	t := &Theme{
		Theme: th,
	}

	t.Theme.TextSize = unit.Sp(DefaultTextSize)
	t.Theme.Palette.Bg = color.NRGBA{R: 37, G: 37, B: 38, A: 255}

	return t
}
