package analyzer

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func TestStructureAnalyzer(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/images/hero.png", Type: asset.TypeImage},
		{Path: "textures/hero.png", Type: asset.TypeImage},
		{Path: "assets/audio/blast.wav", Type: asset.TypeAudio},
		{Path: "sounds/blast.wav", Type: asset.TypeAudio},
	}
	a := &StructureAnalyzer{}
	diags := a.Analyze(assets, &pol)
	if len(diags.Items) != 2 {
		t.Errorf("got %d diagnostics, want 2", len(diags.Items))
		for _, d := range diags.Items {
			t.Logf("  %s", d)
		}
	}
}
