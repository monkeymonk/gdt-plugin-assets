package rename

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manifest struct {
	Timestamp string     `json:"timestamp"`
	Ops       []RenameOp `json:"ops"`
}

func WriteManifest(path string, ops []RenameOp) error {
	m := Manifest{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Ops:       ops,
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadManifest(path string) ([]RenameOp, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m.Ops, nil
}

func Rollback(ops []RenameOp) []error {
	var errs []error
	for _, op := range ops {
		if err := os.MkdirAll(filepath.Dir(op.AbsOld), 0755); err != nil {
			errs = append(errs, fmt.Errorf("mkdir %s: %w", op.AbsOld, err))
			continue
		}
		if err := os.Rename(op.AbsNew, op.AbsOld); err != nil {
			errs = append(errs, fmt.Errorf("rollback %s: %w", op.NewPath, err))
		}
	}
	return errs
}
