package system

import (
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Theme struct {
	*material.Theme
}

func NewTheme() *Theme {
	th := material.NewTheme()
	t := &Theme{
		Theme: th,
	}

	t.Theme.TextSize = unit.Sp(14)
	t.Theme.Palette.Bg = color.NRGBA{R: 37, G: 37, B: 38, A: 255}

	return t
}
