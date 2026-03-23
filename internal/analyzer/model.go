package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

type ModelAnalyzer struct{}

func (a *ModelAnalyzer) Name() string { return "model" }

func (a *ModelAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}

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

		checkOversize(ast, 100*1024*1024, "model.oversize", diags)
	}
	return diags
}
