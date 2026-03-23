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
