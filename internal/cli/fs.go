package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListMarkdownFiles expands the given paths into a list of markdown files.
func ListMarkdownFiles(paths []string) ([]string, error) {
	seen := make(map[string]struct{})
	var files []string

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", p, err)
		}
		if info.IsDir() {
			err = filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if isMarkdown(path) {
					if _, ok := seen[path]; !ok {
						seen[path] = struct{}{}
						files = append(files, path)
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else if isMarkdown(p) {
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				files = append(files, p)
			}
		}
	}

	sort.Strings(files)
	return files, nil
}

func isMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}
