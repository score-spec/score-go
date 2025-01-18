// Copyright 2024 Humanitec
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loader

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

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
				if utf8.Valid(raw) {
					content := string(raw)
					c.Files[i].Content = &content
				} else {
					content := base64.StdEncoding.EncodeToString(raw)
					c.Files[i].BinaryContent = &content
				}
			}
		}
	}

	return nil
}

// readFile reads a text file into memory
func readFile(baseDir, path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
