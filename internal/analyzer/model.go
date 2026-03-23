package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

const modelOversizeFallback = 100 * 1024 * 1024

type ModelAnalyzer struct{}

func (a *ModelAnalyzer) Name() string { return "model" }

func (a *ModelAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}

	maxBytes := int64(modelOversizeFallback)
	if pol.Models.MaxSizeKB > 0 {
		maxBytes = int64(pol.Models.MaxSizeKB) * 1024
	}

	for _, ast := range assets {
		if ast.Type != asset.TypeModel {
			continue
		}
		ext := strings.ToLower(filepath.Ext(ast.Path))

		if pol.Models.WarnOnFBX && ext == ".fbx" {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Warning,
				Rule:        "model.fbx",
				Message:     fmt.Sprintf("FBX format detected; preferred formats: %v", pol.Models.PreferredFormats),
				Explanation: "FBX is proprietary; glTF/GLB preferred for open toolchains",
			})
		}

		checkOversize(ast, maxBytes, "model.oversize", "optimization", diags)
	}
	return diags
}
