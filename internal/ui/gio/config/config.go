package config

import (
	"image/color"
	"time"

	"gioui.org/unit"
)

type UI struct {
	Theme    Theme
	View     View
	Editor   Editor
	Explorer Explorer
}

type Theme struct {
	TextSize unit.Sp
	Bg       color.NRGBA
	Fg       color.NRGBA
}

type View struct {
	ExplorerPaneWidth int
	PaneDividerWidth  int
	ErrorPopup        ErrorPopup
}

type ErrorPopup struct {
	Duration  time.Duration
	MaxWidth  unit.Dp
	Radius    unit.Dp
	Inset     unit.Dp
	PaddingX  unit.Dp
	PaddingY  unit.Dp
	Bg        color.NRGBA
	TextColor color.NRGBA
}

type Editor struct {
	DefaultMode        uint8
	TabWidth           int
	InsetTop           unit.Dp
	InsetBottom        unit.Dp
	InsetLeft          unit.Dp
	LineNumberGap      unit.Dp
	LineNumberPadRight unit.Dp
	LineNumberDigit    unit.Dp
	LineHeight         unit.Dp
	CursorStrokeWidth  unit.Dp
	TextColor          color.NRGBA
	LineNumberColor    color.NRGBA
	CursorColor        color.NRGBA
}

type Explorer struct {
	NodeGap        unit.Dp
	IndentPerDepth int
	RowHeightDelta int
}

func Default() UI {
	return UI{
		Theme: Theme{
			TextSize: unit.Sp(16),
			Bg:       color.NRGBA{R: 37, G: 37, B: 38, A: 255},
			Fg:       color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		},
		View: View{
			ExplorerPaneWidth: 200,
			PaneDividerWidth:  1,
			ErrorPopup: ErrorPopup{
				Duration:  4 * time.Second,
				MaxWidth:  unit.Dp(320),
				Radius:    unit.Dp(10),
				Inset:     unit.Dp(16),
				PaddingX:  unit.Dp(16),
				PaddingY:  unit.Dp(12),
				Bg:        color.NRGBA{R: 177, G: 55, B: 72, A: 235},
				TextColor: color.NRGBA{R: 255, G: 244, B: 246, A: 255},
			},
		},
		Editor: Editor{
			DefaultMode:        1,
			TabWidth:           4,
			InsetTop:           unit.Dp(8),
			InsetBottom:        unit.Dp(8),
			InsetLeft:          unit.Dp(8),
			LineNumberGap:      unit.Dp(10),
			LineNumberPadRight: unit.Dp(10),
			LineNumberDigit:    unit.Dp(10),
			LineHeight:         unit.Dp(22),
			CursorStrokeWidth:  unit.Dp(1),
			TextColor:          color.NRGBA{R: 212, G: 212, B: 212, A: 255},
			LineNumberColor:    color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			CursorColor:        color.NRGBA{R: 120, G: 200, B: 255, A: 180},
		},
		Explorer: Explorer{
			NodeGap:        unit.Dp(4),
			IndentPerDepth: 12,
			RowHeightDelta: 2,
		},
	}
}
