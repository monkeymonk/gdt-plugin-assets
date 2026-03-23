package rename

import (
	"fmt"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

type Collision struct {
	Target  string
	Sources []string
}

type OperationPlan struct {
	Ops        []RenameOp
	Collisions []Collision
	Skipped    int
}

func BuildPlan(assets []asset.Asset, convention string) *OperationPlan {
	ops := Plan(assets, convention)

	existing := make(map[string]bool, len(assets))
	for _, a := range assets {
		existing[a.Path] = true
	}

	renaming := make(map[string]bool, len(ops))
	for _, op := range ops {
		renaming[op.OldPath] = true
	}

	byTarget := make(map[string][]string)
	for _, op := range ops {
		byTarget[op.NewPath] = append(byTarget[op.NewPath], op.OldPath)
	}

	var collisions []Collision

	for target, sources := range byTarget {
		if len(sources) > 1 {
			collisions = append(collisions, Collision{
				Target:  target,
				Sources: sources,
			})
		}
	}

	for _, op := range ops {
		if existing[op.NewPath] && !renaming[op.NewPath] {
			collisions = append(collisions, Collision{
				Target:  op.NewPath,
				Sources: []string{fmt.Sprintf("%s (target already exists)", op.OldPath)},
			})
		}
	}

	return &OperationPlan{
		Ops:        ops,
		Collisions: collisions,
	}
}
