package asset

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func ExtractImageMeta(path string) (*ImageMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return &ImageMeta{
		Width:        cfg.Width,
		Height:       cfg.Height,
		IsPowerOfTwo: isPOT(cfg.Width) && isPOT(cfg.Height),
	}, nil
}

func isPOT(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}
