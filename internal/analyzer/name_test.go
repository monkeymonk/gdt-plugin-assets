package analyzer

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func TestNameAnalyzer_SnakeCase(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/hero_sprite.png", Type: asset.TypeImage},
		{Path: "assets/HeroSprite.png", Type: asset.TypeImage},
		{Path: "assets/hero-sprite.png", Type: asset.TypeImage},
		{Path: "assets/hero sprite.png", Type: asset.TypeImage},
	}
	a := &NameAnalyzer{}
	diags := a.Analyze(assets, &pol)
	if len(diags.Items) != 3 {
		t.Errorf("got %d diagnostics, want 3", len(diags.Items))
		for _, d := range diags.Items {
			t.Logf("  %s", d)
		}
	}
}

func TestNameAnalyzer_SpacesAllowed(t *testing.T) {
	pol := policy.Default()
	pol.Naming.AllowSpaces = true
	assets := []asset.Asset{
		{Path: "assets/hero sprite.png", Type: asset.TypeImage},
	}
	a := &NameAnalyzer{}
	diags := a.Analyze(assets, &pol)
	// allow_spaces=true skips the space check, but "hero sprite" still fails
	// the snake_case regex (spaces not in [a-z0-9_]), so 1 case violation expected.
	if len(diags.Items) != 1 {
		t.Errorf("got %d diagnostics, want 1", len(diags.Items))
	}
}
