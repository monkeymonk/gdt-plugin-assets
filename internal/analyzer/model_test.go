package analyzer

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func TestModelAnalyzer_FBXWarning(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/models/crate.glb", Type: asset.TypeModel},
		{Path: "assets/models/hero.fbx", Type: asset.TypeModel},
	}
	a := &ModelAnalyzer{}
	diags := a.Analyze(assets, &pol)
	if len(diags.Items) != 1 {
		t.Errorf("got %d diagnostics, want 1 (fbx warning)", len(diags.Items))
	}
}
