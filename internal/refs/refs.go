package refs

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var resPathRe = regexp.MustCompile(`path\s*=\s*"(res://[^"]+)"`)

type BrokenRef struct {
	Source string
	Target string
	Line   int
}

// FindBroken walks the project tree to discover engine resource files and checks
// for broken res:// references. Use FindBrokenInFiles when file paths are already known.
func FindBroken(root string) ([]BrokenRef, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".tscn" || ext == ".tres" || ext == ".godot" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return FindBrokenInFiles(root, files), nil
}

// FindBrokenInFiles checks the given file paths for broken res:// references.
// Use this when file paths are already known (e.g. from a prior scan) to avoid
// re-walking the filesystem.
func FindBrokenInFiles(root string, absPaths []string) []BrokenRef {
	var broken []BrokenRef
	for _, path := range absPaths {
		refs, err := scanFileRefs(path, root)
		if err != nil {
			continue
		}
		broken = append(broken, refs...)
	}
	return broken
}

func scanFileRefs(path, root string) ([]BrokenRef, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var broken []BrokenRef
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := resPathRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			resPath := m[1]
			fsPath := resToFS(resPath, root)
			if _, err := os.Stat(fsPath); os.IsNotExist(err) {
				rel, _ := filepath.Rel(root, path)
				broken = append(broken, BrokenRef{
					Source: rel,
					Target: resPath,
					Line:   lineNum,
				})
			}
		}
	}
	return broken, scanner.Err()
}

func resToFS(resPath, root string) string {
	trimmed := strings.TrimPrefix(resPath, "res://")
	return filepath.Join(root, filepath.FromSlash(trimmed))
}
