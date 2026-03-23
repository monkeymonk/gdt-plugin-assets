package scanner

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"assets/images/hero.png":   "fake-png",
		"assets/audio/blast.wav":   "fake-wav",
		"assets/models/crate.glb":  "fake-glb",
		"source_assets/hero.blend": "fake-blend",
		".godot/imported/x.import": "skip",
		".git/objects/abc":         "skip",
		"project.godot":            "[gd_project]",
		"readme.txt":               "hello",
	}
	for rel, content := range files {
		abs := filepath.Join(dir, rel)
		os.MkdirAll(filepath.Dir(abs), 0755)
		os.WriteFile(abs, []byte(content), 0644)
	}
	return dir
}

func TestScanFindsAssets(t *testing.T) {
	dir := setupTestDir(t)
	assets, err := Scan(dir, Options{})
	if err != nil {
		t.Fatal(err)
	}
	// Should find: hero.png, blast.wav, crate.glb, hero.blend
	// Should skip: .godot/*, .git/*, readme.txt
	// Note: project.godot has .godot extension which maps to TypeEngineResource
	// So we expect 5 assets total (hero.png, blast.wav, crate.glb, hero.blend, project.godot)
	if len(assets) != 5 {
		t.Errorf("got %d assets, want 5", len(assets))
		for _, a := range assets {
			t.Logf("  %s (%s)", a.Path, a.Type)
		}
	}
}

func TestScanWithTypeFilter(t *testing.T) {
	dir := setupTestDir(t)
	assets, err := Scan(dir, Options{Types: []string{"image"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 1 {
		t.Errorf("got %d image assets, want 1", len(assets))
	}
}

func TestScan_ImageMetadata(t *testing.T) {
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "assets", "images")
	os.MkdirAll(imgDir, 0755)

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	f, _ := os.Create(filepath.Join(imgDir, "icon.png"))
	png.Encode(f, img)
	f.Close()

	assets, err := Scan(dir, Options{})
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, a := range assets {
		if a.Image != nil && a.Image.Width == 64 && a.Image.Height == 64 {
			found = true
			if !a.Image.IsPowerOfTwo {
				t.Error("64x64 should be POT")
			}
		}
	}
	if !found {
		t.Error("image metadata not populated for icon.png")
	}
}

func TestScanWithHash(t *testing.T) {
	dir := setupTestDir(t)
	assets, err := Scan(dir, Options{Hash: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range assets {
		if a.Hash == "" {
			t.Errorf("asset %s missing hash", a.Path)
		}
	}
}
