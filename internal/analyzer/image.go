package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

const imageOversizeThreshold = 10 * 1024 * 1024

type ImageAnalyzer struct{}

func (a *ImageAnalyzer) Name() string { return "image" }

func (a *ImageAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}
	allowedSet := buildFormatSet(pol.Images.AllowedFormats)

	for _, ast := range assets {
		if ast.Type != asset.TypeImage {
			continue
		}
		ext := strings.ToLower(filepath.Ext(ast.Path))

		if len(allowedSet) > 0 && !allowedSet[ext] {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Warning,
				Rule:        "image.format",
				Message:     fmt.Sprintf("format %s not in allowed list %v", ext, pol.Images.AllowedFormats),
				Explanation: "Policy restricts image formats for consistency",
			})
		}

		checkOversize(ast, imageOversizeThreshold, "image.oversize", diags)
	}
	return diags
}
