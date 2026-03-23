package refs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepairPlan(t *testing.T) {
	dir := t.TempDir()

	scene := "[gd_scene format=3]\n" +
		"[ext_resource type=\"Texture2D\" path=\"res://assets/old_name.png\" id=\"1\"]\n" +
		"[ext_resource type=\"PackedScene\" path=\"res://scenes/valid.tscn\" id=\"2\"]\n"
	scenePath := filepath.Join(dir, "test.tscn")
	os.WriteFile(scenePath, []byte(scene), 0644)

	mapping := map[string]string{
		"res://assets/old_name.png": "res://assets/new_name.png",
	}

	ops, err := PlanRepair(scenePath, mapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 repair op, got %d", len(ops))
	}
	if ops[0].OldRef != "res://assets/old_name.png" {
		t.Errorf("unexpected old ref: %s", ops[0].OldRef)
	}
	if ops[0].NewRef != "res://assets/new_name.png" {
		t.Errorf("unexpected new ref: %s", ops[0].NewRef)
	}
}

func TestApplyRepair(t *testing.T) {
	dir := t.TempDir()

	scene := "[ext_resource type=\"Texture2D\" path=\"res://assets/old.png\" id=\"1\"]"
	scenePath := filepath.Join(dir, "test.tscn")
	os.WriteFile(scenePath, []byte(scene), 0644)

	ops := []RepairOp{
		{File: scenePath, Line: 1, OldRef: "res://assets/old.png", NewRef: "res://assets/new.png"},
	}

	if err := ApplyRepair(ops); err != nil {
		t.Fatalf("apply error: %v", err)
	}

	data, _ := os.ReadFile(scenePath)
	if !strings.Contains(string(data), "res://assets/new.png") {
		t.Error("expected ref to be rewritten")
	}
	if strings.Contains(string(data), "res://assets/old.png") {
		t.Error("old ref should be gone")
	}
}
