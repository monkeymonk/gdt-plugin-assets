package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

type AudioAnalyzer struct{}

func (a *AudioAnalyzer) Name() string { return "audio" }

func (a *AudioAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}
	preferred := buildFormatSet(pol.Audio.PreferredFormats)

	for _, ast := range assets {
		if ast.Type != asset.TypeAudio {
			continue
		}
		ext := strings.ToLower(filepath.Ext(ast.Path))
		if len(preferred) > 0 && !preferred[ext] {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Warning,
				Rule:        "audio.format",
				Message:     fmt.Sprintf("format %s not in preferred list %v", ext, pol.Audio.PreferredFormats),
				Explanation: "Non-preferred audio formats may cause runtime issues",
			})
		}

		checkOversize(ast, 50*1024*1024, "audio.oversize", diags)
	}
	return diags
}
