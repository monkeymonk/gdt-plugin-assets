package analyzer

import (
	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
	"github.com/monkeymonk/gdt-assets/internal/policy"
)

type Analyzer interface {
	Name() string
	Analyze(assets []asset.Asset, pol *policy.Policy) *diagnostic.Set
}

func RunAll(analyzers []Analyzer, assets []asset.Asset, pol *policy.Policy) *diagnostic.Set {
	result := &diagnostic.Set{}
	for _, a := range analyzers {
		result.AddAll(a.Analyze(assets, pol))
	}
	return result
}

func DefaultAnalyzers() []Analyzer {
	return []Analyzer{
		&NameAnalyzer{},
		&StructureAnalyzer{},
		&ImageAnalyzer{},
		&AudioAnalyzer{},
		&ModelAnalyzer{},
	}
}
