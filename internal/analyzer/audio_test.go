package analyzer

import (
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

func TestAudioAnalyzer_FormatWarning(t *testing.T) {
	pol := policy.Default()
	assets := []asset.Asset{
		{Path: "assets/audio/music.mp3", Type: asset.TypeAudio, Size: 1024},
		{Path: "assets/audio/sfx.ogg", Type: asset.TypeAudio, Size: 1024},
	}
	a := &AudioAnalyzer{}
	diags := a.Analyze(assets, &pol)
	if len(diags.Items) != 1 {
		t.Errorf("got %d diagnostics, want 1", len(diags.Items))
	}
}
