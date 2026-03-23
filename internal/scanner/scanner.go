package scanner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/monkeymonk/gdt-assets/internal/asset"
)

var skipDirs = map[string]bool{
	".godot": true, ".git": true, ".import": true,
	"__pycache__": true, "node_modules": true,
}

type Options struct {
	Types []string
	Hash  bool
}

func Scan(root string, opts Options) ([]asset.Asset, error) {
	typeFilter := make(map[asset.AssetType]bool)
	for _, t := range opts.Types {
		for at := asset.AssetType(1); at <= asset.TypeEngineResource; at++ {
			if at.String() == t {
				typeFilter[at] = true
			}
		}
	}

	var assets []asset.Asset
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		at := asset.DetectType(path)
		if at == asset.TypeUnknown {
			return nil
		}
		if len(typeFilter) > 0 && !typeFilter[at] {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "" {
			rel = path
		}
		rel = filepath.ToSlash(rel)

		a := asset.Asset{
			Path:    rel,
			AbsPath: path,
			Type:    at,
			Size:    info.Size(),
		}

		if opts.Hash {
			h, herr := hashFile(path)
			if herr == nil {
				a.Hash = h
			}
		}

		if a.Type == asset.TypeImage {
			meta, err := asset.ExtractImageMeta(path)
			if err == nil {
				a.Image = meta
			}
		}

		assets = append(assets, a)
		return nil
	})
	return assets, err
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
