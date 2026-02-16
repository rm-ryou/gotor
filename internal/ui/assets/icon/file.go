package icon

import (
	"image/color"
)

const (
	ArrowCollapsed = ""
	ArrowExpanded  = ""
)

type nerdIconStyle struct {
	Glyph string
	Color color.NRGBA
}

var (
	DefaultFileIcon = nerdIconStyle{
		Glyph: "󰈔",
		Color: color.NRGBA{R: 180, G: 180, B: 180, A: 255},
	}
	FolderClosedIcon = nerdIconStyle{
		Glyph: "",
		Color: color.NRGBA{R: 224, G: 187, B: 93, A: 255},
	}
	FolderOpenIcon = nerdIconStyle{
		Glyph: "",
		Color: color.NRGBA{R: 224, G: 187, B: 93, A: 255},
	}
)
