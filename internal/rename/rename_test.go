package rename

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"HeroSprite", "hero_sprite"},
		{"hero-sprite", "hero_sprite"},
		{"hero sprite", "hero_sprite"},
		{"hero_sprite", "hero_sprite"},
		{"MyFBXFile", "my_fbx_file"},
	}
	for _, tt := range tests {
		got := ToSnakeCase(tt.in)
		if got != tt.want {
			t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestPlanRenames(t *testing.T) {
	assets := []asset.Asset{
		{Path: "assets/HeroSprite.png", AbsPath: "/p/assets/HeroSprite.png"},
		{Path: "assets/hero_icon.png", AbsPath: "/p/assets/hero_icon.png"},
	}
	plan := Plan(assets, "snake")
	if len(plan) != 1 {
		t.Errorf("got %d renames, want 1", len(plan))
	}
	if len(plan) > 0 && plan[0].NewPath != "assets/hero_sprite.png" {
		t.Errorf("new path = %q, want assets/hero_sprite.png", plan[0].NewPath)
	}
}
