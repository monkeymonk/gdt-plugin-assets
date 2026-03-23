package analyzer

import (
	"fmt"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
)

func buildFormatSet(formats []string) map[string]bool {
	m := make(map[string]bool, len(formats))
	for _, f := range formats {
		m["."+strings.ToLower(f)] = true
	}
	return m
}

func checkOversize(a asset.Asset, threshold int64, rule, category string, diags *diagnostic.Set) {
	if a.Size >= threshold {
		diags.Add(diagnostic.Diagnostic{
			Path:     a.Path,
			Severity: diagnostic.Error,
			Rule:     rule,
			Category: category,
			Message:  fmt.Sprintf("file is %s, exceeds %s threshold", asset.HumanSize(a.Size), asset.HumanSize(threshold)),
		})
	}
}
