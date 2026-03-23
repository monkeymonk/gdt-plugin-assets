package analyzer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

var (
	snakeRe  = regexp.MustCompile(`^[a-z0-9][a-z0-9_]*$`)
	kebabRe  = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]*$`)
	camelRe  = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
	pascalRe = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
)

type NameAnalyzer struct{}

func (a *NameAnalyzer) Name() string { return "naming" }

func (a *NameAnalyzer) Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	diags := &diagnostic.Set{}
	caseCheck := caseChecker(pol.Naming.Case)

	for _, ast := range assets {
		base := filepath.Base(ast.Path)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		if !pol.Naming.AllowSpaces && strings.Contains(name, " ") {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Error,
				Rule:        "naming.no_spaces",
				Message:     fmt.Sprintf("filename contains spaces: %q", base),
				Explanation: "Spaces in asset filenames cause path issues across platforms",
				CanAutoFix:  true,
			})
			continue
		}

		if caseCheck != nil && !caseCheck(name) {
			diags.Add(diagnostic.Diagnostic{
				Path:        ast.Path,
				Severity:    diagnostic.Warning,
				Rule:        "naming.case",
				Message:     fmt.Sprintf("filename %q does not match %s convention", base, pol.Naming.Case),
				Explanation: fmt.Sprintf("Project policy requires %s naming", pol.Naming.Case),
				CanAutoFix:  true,
			})
		}
	}
	return diags
}

func caseChecker(convention string) func(string) bool {
	switch convention {
	case "snake":
		return snakeRe.MatchString
	case "kebab":
		return kebabRe.MatchString
	case "camel":
		return camelRe.MatchString
	case "pascal":
		return pascalRe.MatchString
	default:
		return nil
	}
}
