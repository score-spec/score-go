package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/score-spec/score-go/types"
)

// Normalize normalizes the target Workload by:
// * embedding container file sources as content
func Normalize(w *types.Workload, baseDir string) error {
	for name, c := range w.Containers {
		for i, f := range c.Files {
			if f.Source != nil {
				raw, err := readFile(baseDir, *f.Source)
				if err != nil {
					return fmt.Errorf("embedding file '%s' for container '%s': %w", *f.Source, name, err)
				}

				c.Files[i].Source = nil
				c.Files[i].Content = &raw
			}
		}
	}

	return nil
}

// readFile reads a text file into memory
func readFile(baseDir, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
