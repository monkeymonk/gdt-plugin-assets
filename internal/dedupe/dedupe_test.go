package dedupe

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

func TestFindExactDuplicates(t *testing.T) {
	assets := []asset.Asset{
		{Path: "a/hero.png", Hash: "abc123"},
		{Path: "b/hero_copy.png", Hash: "abc123"},
		{Path: "c/other.png", Hash: "def456"},
	}
	groups := FindExact(assets)
	if len(groups) != 1 {
		t.Errorf("got %d groups, want 1", len(groups))
	}
	if len(groups) > 0 && len(groups[0].Paths) != 2 {
		t.Errorf("group has %d items, want 2", len(groups[0].Paths))
	}
}

func TestFindNameDuplicates(t *testing.T) {
	assets := []asset.Asset{
		{Path: "a/hero.png"},
		{Path: "b/hero.png"},
		{Path: "c/villain.png"},
	}
	groups := FindByName(assets)
	if len(groups) != 1 {
		t.Errorf("got %d groups, want 1", len(groups))
	}
}
