package analyzer

import (
	"fmt"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

type StructureAnalyzer struct{}

func (a *StructureAnalyzer) Name() string { return "structure" }

func (a *StructureAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}
	expected := map[asset.AssetType]string{
		asset.TypeImage:  pol.Folders.Images,
		asset.TypeAudio:  pol.Folders.Audio,
		asset.TypeModel:  pol.Folders.Models,
		asset.TypeVector: pol.Folders.Vectors,
		asset.TypeFont:   pol.Folders.Fonts,
	}

	for _, ast := range assets {
		folder, ok := expected[ast.Type]
		if !ok || folder == "" {
			continue
		}
		if !strings.HasPrefix(ast.Path, folder+"/") && !strings.HasPrefix(ast.Path, folder+"\\") {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Warning,
				Rule:        "structure.folder",
				Message:     fmt.Sprintf("expected under %s/, found at %s", folder, ast.Path),
				Explanation: fmt.Sprintf("Policy expects %s assets in %s/", ast.Type, folder),
				CanAutoFix:  true,
			})
		}
	}
	return diags
}
