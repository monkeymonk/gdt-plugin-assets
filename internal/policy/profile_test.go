package policy

import "testing"

func TestResolveProfile_OverridesImageMaxSize(t *testing.T) {
	pol := Default()
	pol.Profiles = map[string]ProfileOverrides{
		"mobile": {
			Images: &ImageOverrides{MaxSizeDefaultKB: intPtr(1024)},
		},
	}

	resolved, err := ResolveProfile(&pol, "mobile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.Images.MaxSizeDefaultKB != 1024 {
		t.Errorf("expected MaxSizeDefaultKB=1024, got %d", resolved.Images.MaxSizeDefaultKB)
	}
	if resolved.Images.RequirePowerOfTwo != true {
		t.Error("expected RequirePowerOfTwo to remain true")
	}
}

func TestResolveProfile_UnknownProfile(t *testing.T) {
	pol := Default()
	_, err := ResolveProfile(&pol, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown profile")
	}
}

func TestResolveProfile_EmptyName(t *testing.T) {
	pol := Default()
	resolved, err := ResolveProfile(&pol, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.Images.MaxSizeDefaultKB != pol.Images.MaxSizeDefaultKB {
		t.Error("empty profile should return base policy unchanged")
	}
}

func intPtr(n int) *int { return &n }
