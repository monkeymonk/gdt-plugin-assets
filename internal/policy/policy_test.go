package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPolicy(t *testing.T) {
	p := Default()
	if p.Naming.Case != "snake" {
		t.Errorf("default naming case = %q, want snake", p.Naming.Case)
	}
	if p.Images.MaxSizeDefaultKB != 4096 {
		t.Errorf("default image max = %d, want 4096", p.Images.MaxSizeDefaultKB)
	}
}

func TestLoadPolicy(t *testing.T) {
	dir := t.TempDir()
	content := `
version = 1

[naming]
case = "kebab"
allow_spaces = true

[images]
max_size_default_kb = 2048
`
	path := filepath.Join(dir, "assets.policy.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	p, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if p.Naming.Case != "kebab" {
		t.Errorf("naming case = %q, want kebab", p.Naming.Case)
	}
	if !p.Naming.AllowSpaces {
		t.Error("allow_spaces should be true")
	}
	if p.Images.MaxSizeDefaultKB != 2048 {
		t.Errorf("image max = %d, want 2048", p.Images.MaxSizeDefaultKB)
	}
	// Defaults should fill unset fields
	if p.Images.MaxSizeUIKB != 2048 {
		t.Errorf("image max_ui should default to 2048, got %d", p.Images.MaxSizeUIKB)
	}
}

func TestLoadOrDefault(t *testing.T) {
	p := LoadOrDefault("/nonexistent/path")
	if p.Version != 1 {
		t.Errorf("version = %d, want 1", p.Version)
	}
}
