package assets

import (
	_ "embed"
)

var (
	//go:embed blocks.png
	BlocksPng []byte

	//go:embed logo.png
	LogoPng []byte
)
