package asset

import "testing"

func TestDetectType(t *testing.T) {
	tests := []struct {
		path string
		want AssetType
	}{
		{"hero.png", TypeImage},
		{"icon.svg", TypeVector},
		{"blast.wav", TypeAudio},
		{"cutscene.mp4", TypeVideo},
		{"crate.glb", TypeModel},
		{"hero_idle.anim", TypeAnimation},
		{"ui.ttf", TypeFont},
		{"fx.gdshader", TypeShader},
		{"hero.blend", TypeDocument},
		{"level.tscn", TypeEngineResource},
		{"readme.txt", TypeUnknown},
	}
	for _, tt := range tests {
		got := DetectType(tt.path)
		if got != tt.want {
			t.Errorf("DetectType(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}
