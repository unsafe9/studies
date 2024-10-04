package main

import (
	"io/fs"
	"iter"
	"path/filepath"
)

func glob(patterns []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		visited := make(map[string]struct{})
		for _, pattern := range patterns {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return
			}
			for _, match := range matches {
				match = filepath.Clean(match)
				if _, ok := visited[match]; ok {
					continue
				}
				visited[match] = struct{}{}

				if !yield(match) {
					return
				}
			}
		}
	}
}

func walkFiles(dirs []string, files []string, exts []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		visited := make(map[string]struct{})
		stop := false

		for _, dir := range dirs {
			err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if len(exts) > 0 {
					ext := filepath.Ext(path)
					if len(ext) > 0 && !slices.Contains(exts, ext[1:]) {
						return nil
					}
				}

				path, err = filepath.Rel(".", path)
				if err != nil {
					return nil
				}

				if _, ok := visited[path]; ok {
					return nil
				}
				visited[path] = struct{}{}

				if !yield(path) {
					stop = true
					return filepath.SkipAll
				}

				return nil
			})
			if err != nil {
				return
			}
			if stop {
				return
			}
		}

		for _, path := range files {
			var err error
			path, err = filepath.Rel(".", path)
			if err != nil {
				continue
			}

			if _, ok := visited[path]; ok {
				continue
			}
			visited[path] = struct{}{}

			if !yield(path) {
				return
			}
		}
	}
}
