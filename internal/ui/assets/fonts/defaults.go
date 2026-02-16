package fonts

import _ "embed"

const DefaultTypeface = "Hack"

var (
	//go:embed Hack/HackNerdFontMono-Regular.ttf
	hackNerdFontMonoRegular []byte
	//go:embed Hack/HackNerdFontMono-Bold.ttf
	hackNerdFontMonoBold []byte
)

func DefaultHackMonoRegular() []byte {
	return append([]byte(nil), hackNerdFontMonoRegular...)
}

func DefaultHackMonoBold() []byte {
	return append([]byte(nil), hackNerdFontMonoBold...)
}
