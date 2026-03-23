package analyzer

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func TestImageAnalyzer_OversizedFile(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/images/small.png", Type: asset.TypeImage, Size: 100 * 1024},
		{Path: "assets/images/huge.png", Type: asset.TypeImage, Size: 50 * 1024 * 1024},
	}
	a := &ImageAnalyzer{}
	diags := a.Analyze(assets, &pol)
	if len(diags.Items) < 1 {
		t.Error("expected at least 1 diagnostic for oversized image")
	}
}

func TestImageAnalyzer_PolicyMaxSize(t *testing.T) {
	pol := policy.Default()
	pol.Images.MaxSizeDefaultKB = 100 // 100 KB

	assets := []asset.Asset{
		{Path: "assets/images/big.png", Type: asset.TypeImage, Size: 200 * 1024}, // 200 KB
	}

	diags := (&ImageAnalyzer{}).Analyze(assets, &pol)

	found := false
	for _, d := range diags.Items {
		if d.Rule == "image.oversize" {
			found = true
		}
	}
	if !found {
		t.Error("expected image.oversize diagnostic for 200KB file with 100KB policy limit")
	}
}

func TestImageAnalyzer_PowerOfTwo(t *testing.T) {
	pol := policy.Default()
	pol.Images.RequirePowerOfTwo = true
	pol.Images.AllowNonPotForUI = true

	assets := []asset.Asset{
		{
			Path: "assets/images/texture.png", Type: asset.TypeImage, Size: 1024,
			Image: &asset.ImageMeta{Width: 300, Height: 300, IsPowerOfTwo: false},
		},
		{
			Path: "assets/images/ui/button.png", Type: asset.TypeImage, Size: 1024,
			Image: &asset.ImageMeta{Width: 300, Height: 200, IsPowerOfTwo: false},
		},
		{
			Path: "assets/images/tile.png", Type: asset.TypeImage, Size: 1024,
			Image: &asset.ImageMeta{Width: 512, Height: 512, IsPowerOfTwo: true},
		},
	}

	diags := (&ImageAnalyzer{}).Analyze(assets, &pol)

	potCount := 0
	for _, d := range diags.Items {
		if d.Rule == "image.pot" {
			potCount++
		}
	}
	// texture.png should fail POT, ui/button.png should be exempt, tile.png is POT
	if potCount != 1 {
		t.Errorf("expected 1 POT diagnostic, got %d", potCount)
	}
}

func TestImageAnalyzer_DisallowedFormat(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/images/photo.bmp", Type: asset.TypeImage, Size: 1024},
	}
	a := &ImageAnalyzer{}
	diags := a.Analyze(assets, &pol)
	found := false
	for _, d := range diags.Items {
		if d.Rule == "image.format" {
			found = true
		}
	}
	if !found {
		t.Error("expected format violation for .bmp")
	}
}
