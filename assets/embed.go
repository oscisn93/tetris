package assets

import (
	_ "embed"
)

var (
	//go:embed blocks.png
	BlocksPng []byte

	//go:embed logo.png
	LogoPng []byte

	//go:embed ebitengine.png
	EbitenginePng []byte

	//go:embed pressstart2p.ttf
	PerssStart2PTtf []byte
)
