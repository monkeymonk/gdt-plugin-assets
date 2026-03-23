package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

const imageOversizeFallback = 10 * 1024 * 1024 // 10 MB fallback

type ImageAnalyzer struct{}

func (a *ImageAnalyzer) Name() string { return "image" }

func (a *ImageAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}
	allowedSet := buildFormatSet(pol.Images.AllowedFormats)

	maxBytes := int64(imageOversizeFallback)
	if pol.Images.MaxSizeDefaultKB > 0 {
		maxBytes = int64(pol.Images.MaxSizeDefaultKB) * 1024
	}

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

		checkOversize(ast, maxBytes, "image.oversize", diags)

		if pol.Images.RequirePowerOfTwo && ast.Image != nil && !ast.Image.IsPowerOfTwo {
			if pol.Images.AllowNonPotForUI && isUIPath(ast.Path) {
				// exempt UI assets
			} else {
				diags.Add(diagnostic.Diagnostic{
					Path:        ast.Path,
					Severity:    diagnostic.Warning,
					Rule:        "image.pot",
					Message:     fmt.Sprintf("dimensions %dx%d are not power-of-two", ast.Image.Width, ast.Image.Height),
					Explanation: "Non-POT textures may cause GPU memory waste or compatibility issues",
				})
			}
		}
	}
	return diags
}

func isUIPath(path string) bool {
	return strings.Contains(strings.ToLower(path), "/ui/") ||
		strings.Contains(strings.ToLower(path), "/gui/") ||
		strings.Contains(strings.ToLower(path), "/hud/")
}
