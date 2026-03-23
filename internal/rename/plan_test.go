package rename

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

func TestBuildPlan_NoCollisions(t *testing.T) {
	assets := []asset.Asset{
		{Path: "assets/MyTexture.png", AbsPath: "/proj/assets/MyTexture.png", Type: asset.TypeImage},
		{Path: "assets/OtherFile.png", AbsPath: "/proj/assets/OtherFile.png", Type: asset.TypeImage},
	}

	plan := BuildPlan(assets, "snake")
	if len(plan.Collisions) != 0 {
		t.Errorf("expected 0 collisions, got %d", len(plan.Collisions))
	}
	if len(plan.Ops) != 2 {
		t.Errorf("expected 2 ops, got %d", len(plan.Ops))
	}
}

func TestBuildPlan_DetectsCollision(t *testing.T) {
	assets := []asset.Asset{
		{Path: "assets/my-file.png", AbsPath: "/proj/assets/my-file.png", Type: asset.TypeImage},
		{Path: "assets/My File.png", AbsPath: "/proj/assets/My File.png", Type: asset.TypeImage},
	}

	plan := BuildPlan(assets, "snake")
	if len(plan.Collisions) == 0 {
		t.Error("expected collision detection for same target")
	}
}

func TestBuildPlan_DetectsExistingTarget(t *testing.T) {
	assets := []asset.Asset{
		{Path: "assets/My File.png", AbsPath: "/proj/assets/My File.png", Type: asset.TypeImage},
		{Path: "assets/my_file.png", AbsPath: "/proj/assets/my_file.png", Type: asset.TypeImage},
	}

	plan := BuildPlan(assets, "snake")
	if len(plan.Collisions) == 0 {
		t.Error("expected collision: target already exists as different source")
	}
}
