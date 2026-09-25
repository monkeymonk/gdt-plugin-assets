package refs

import (
	"bufio"
	"os"
	"strings"
)

type RepairOp struct {
	File   string
	Line   int
	OldRef string
	NewRef string
}

type RenamePair struct {
	OldPath string
	NewPath string
}

func PlanRepair(path string, mapping map[string]string) ([]RepairOp, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ops []RepairOp
	sc := bufio.NewScanner(f)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		matches := resPathRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			oldRef := m[1]
			if newRef, ok := mapping[oldRef]; ok {
				ops = append(ops, RepairOp{
					File:   path,
					Line:   lineNum,
					OldRef: oldRef,
					NewRef: newRef,
				})
			}
		}
	}
	return ops, sc.Err()
}

// PlanRepairFromRenames plans repairs for all files, given a set of
// renamed paths. Returns the planned ops plus any per-file errors
// encountered (e.g. an unreadable .tscn/.tres file) so the caller can
// warn the user rather than silently repairing fewer references than
// actually exist — a file that fails to plan is not the same as a
// file with nothing to repair.
func PlanRepairFromRenames(root string, renames []RenamePair, files []string) ([]RepairOp, map[string]error, error) {
	mapping := make(map[string]string, len(renames))
	for _, r := range renames {
		mapping["res://"+r.OldPath] = "res://" + r.NewPath
	}

	var allOps []RepairOp
	var skipped map[string]error
	for _, f := range files {
		ops, err := PlanRepair(f, mapping)
		if err != nil {
			if skipped == nil {
				skipped = make(map[string]error)
			}
			skipped[f] = err
			continue
		}
		allOps = append(allOps, ops...)
	}
	return allOps, skipped, nil
}

func ApplyRepair(ops []RepairOp) error {
	byFile := make(map[string][]RepairOp)
	for _, op := range ops {
		byFile[op.File] = append(byFile[op.File], op)
	}

	for path, fileOps := range byFile {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		for _, op := range fileOps {
			old := `path="` + op.OldRef + `"`
			new := `path="` + op.NewRef + `"`
			content = strings.ReplaceAll(content, old, new)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}
