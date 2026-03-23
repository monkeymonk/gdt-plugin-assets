package dedupe

import (
	"path/filepath"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

type DupeGroup struct {
	Key   string
	Paths []string
}

func FindExact(assets []asset.Asset) []DupeGroup {
	hasAnyHash := false
	for _, a := range assets {
		if a.Hash != "" {
			hasAnyHash = true
			break
		}
	}
	if !hasAnyHash {
		return nil
	}

	byHash := make(map[string][]string)
	for _, a := range assets {
		if a.Hash == "" {
			continue
		}
		byHash[a.Hash] = append(byHash[a.Hash], a.Path)
	}
	var groups []DupeGroup
	for hash, paths := range byHash {
		if len(paths) > 1 {
			groups = append(groups, DupeGroup{Key: hash, Paths: paths})
		}
	}
	return groups
}

func FindByName(assets []asset.Asset) []DupeGroup {
	byName := make(map[string][]string)
	for _, a := range assets {
		name := strings.ToLower(filepath.Base(a.Path))
		byName[name] = append(byName[name], a.Path)
	}
	var groups []DupeGroup
	for name, paths := range byName {
		if len(paths) > 1 {
			groups = append(groups, DupeGroup{Key: name, Paths: paths})
		}
	}
	return groups
}
