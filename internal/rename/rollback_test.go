package rename

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndLoadManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rollback.json")

	ops := []RenameOp{
		{OldPath: "a/Old.png", NewPath: "a/old.png", AbsOld: "/p/a/Old.png", AbsNew: "/p/a/old.png"},
	}

	if err := WriteManifest(path, ops); err != nil {
		t.Fatalf("write: %v", err)
	}

	loaded, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 op, got %d", len(loaded))
	}
	if loaded[0].OldPath != "a/Old.png" || loaded[0].NewPath != "a/old.png" {
		t.Errorf("unexpected op: %+v", loaded[0])
	}
}

func TestRollback(t *testing.T) {
	dir := t.TempDir()

	newDir := filepath.Join(dir, "a")
	os.MkdirAll(newDir, 0755)
	newPath := filepath.Join(newDir, "old.png")
	os.WriteFile(newPath, []byte("data"), 0644)

	ops := []RenameOp{
		{
			OldPath: "a/Old.png", NewPath: "a/old.png",
			AbsOld: filepath.Join(newDir, "Old.png"),
			AbsNew: newPath,
		},
	}

	errs := Rollback(ops)
	if len(errs) != 0 {
		t.Fatalf("rollback errors: %v", errs)
	}

	if _, err := os.Stat(filepath.Join(newDir, "Old.png")); err != nil {
		t.Error("expected Old.png to be restored")
	}
	if _, err := os.Stat(newPath); err == nil {
		t.Error("expected old.png to be gone after rollback")
	}
}
