package refs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindBrokenRefs(t *testing.T) {
	dir := t.TempDir()
	tscn := `[gd_scene format=3]
[ext_resource type="Texture2D" path="res://assets/hero.png" id="1"]
[ext_resource type="Texture2D" path="res://assets/missing.png" id="2"]
`
	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "level.tscn"), []byte(tscn), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "hero.png"), []byte("x"), 0644)

	broken, err := FindBroken(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(broken) != 1 {
		t.Errorf("got %d broken refs, want 1", len(broken))
	}
	if len(broken) > 0 && broken[0].Target != "res://assets/missing.png" {
		t.Errorf("wrong target: %s", broken[0].Target)
	}
}
