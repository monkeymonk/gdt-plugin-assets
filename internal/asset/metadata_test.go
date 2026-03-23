package asset

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func writePNG(t *testing.T, dir, name string, w, h int) string {
	t.Helper()
	path := filepath.Join(dir, name)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.White)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestExtractImageMeta(t *testing.T) {
	dir := t.TempDir()
	path := writePNG(t, dir, "test.png", 300, 128)

	meta, err := ExtractImageMeta(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Width != 300 {
		t.Errorf("width = %d, want 300", meta.Width)
	}
	if meta.Height != 128 {
		t.Errorf("height = %d, want 128", meta.Height)
	}
	if meta.IsPowerOfTwo {
		t.Error("300x128 has non-POT width, want IsPowerOfTwo=false")
	}
}

func TestExtractImageMeta_POT(t *testing.T) {
	dir := t.TempDir()
	path := writePNG(t, dir, "pot.png", 512, 512)

	meta, err := ExtractImageMeta(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !meta.IsPowerOfTwo {
		t.Error("512x512 should be POT")
	}
}

func TestExtractImageMeta_NonImage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not_an_image.txt")
	os.WriteFile(path, []byte("hello"), 0644)

	_, err := ExtractImageMeta(path)
	if err == nil {
		t.Error("expected error for non-image file")
	}
}
