package system

import (
	"fmt"
	"strings"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
	assetfonts "github.com/rm-ryou/gotor/internal/ui/assets/fonts"
)

func Prepare() ([]font.FontFace, error) {
	regularFace, err := parseFont(assetfonts.DefaultHackMonoRegular())
	if err != nil {
		return nil, err
	}

	boldFace, err := parseFont(assetfonts.DefaultHackMonoBold())
	if err != nil {
		return nil, err
	}

	return []font.FontFace{
		{Font: font.Font{Typeface: font.Typeface(assetfonts.DefaultTypeface), Weight: font.Normal}, Face: regularFace},
		{Font: font.Font{Typeface: font.Typeface(assetfonts.DefaultTypeface), Weight: font.Bold}, Face: boldFace},
	}, nil
}

func NewShaper(fontFaces []font.FontFace) *text.Shaper {
	return text.NewShaper(text.WithCollection(fontFaces))
}

func DefaultTypefaceWithFallback() string {
	families := []string{assetfonts.DefaultTypeface}

	families = append(families, "Apple Color Emoji")
	return strings.Join(families, ", ")
}

func parseFont(src []byte) (opentype.Face, error) {
	face, err := opentype.Parse(src)
	if err != nil {
		return opentype.Face{}, fmt.Errorf("failed to parse font: %w", err)
	}
	return face, nil
}
