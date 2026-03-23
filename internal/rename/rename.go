package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

type RenameOp struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	AbsOld  string `json:"abs_old"`
	AbsNew  string `json:"abs_new"`
}

var camelSplitRe = regexp.MustCompile(`([a-z0-9])([A-Z])`)
var camelAcronymRe = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)

func ToSnakeCase(s string) string {
	s = camelAcronymRe.ReplaceAllString(s, "${1}_${2}")
	s = camelSplitRe.ReplaceAllString(s, "${1}_${2}")
	s = strings.NewReplacer("-", "_", " ", "_").Replace(s)
	s = strings.ToLower(s)
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return s
}

func Plan(assets []asset.Asset, convention string) []RenameOp {
	converter := getConverter(convention)
	if converter == nil {
		return nil
	}

	var ops []RenameOp
	for _, a := range assets {
		dir := filepath.Dir(a.Path)
		base := filepath.Base(a.Path)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)

		newName := converter(name)
		if newName == name {
			continue
		}

		newBase := newName + ext
		newPath := filepath.Join(dir, newBase)
		newPath = filepath.ToSlash(newPath)
		ops = append(ops, RenameOp{
			OldPath: a.Path,
			NewPath: newPath,
			AbsOld:  a.AbsPath,
			AbsNew:  filepath.Join(filepath.Dir(a.AbsPath), newBase),
		})
	}
	return ops
}

func Apply(ops []RenameOp) []error {
	var errs []error
	for _, op := range ops {
		if err := os.MkdirAll(filepath.Dir(op.AbsNew), 0755); err != nil {
			errs = append(errs, fmt.Errorf("mkdir %s: %w", op.AbsNew, err))
			continue
		}
		if err := os.Rename(op.AbsOld, op.AbsNew); err != nil {
			errs = append(errs, fmt.Errorf("rename %s: %w", op.OldPath, err))
		}
	}
	return errs
}

func getConverter(convention string) func(string) string {
	switch convention {
	case "snake":
		return ToSnakeCase
	case "kebab":
		return func(s string) string {
			return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
		}
	case "lower":
		return strings.ToLower
	default:
		return nil
	}
}
