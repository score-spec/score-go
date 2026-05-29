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
	"reflect"
	"unicode/utf8"

	"github.com/score-spec/score-go/types"
)

// Normalize normalizes the target Workload by:
// * embedding container file sources as content
func Normalize(w *types.Workload, baseDir string) error {
	for name, c := range w.Containers {
		for target, f := range c.Files {
			updated, changed, err := normalizeContainerFile(f, baseDir)
			if err != nil {
				return fmt.Errorf("embedding file for target '%s' in container '%s': %w", target, name, err)
			}
			if changed {
				c.Files[target] = updated
			}
		}
	}

	return nil
}

func normalizeContainerFile(file any, baseDir string) (any, bool, error) {
	switch f := file.(type) {
	case string:
		raw, err := readFile(baseDir, f)
		if err != nil {
			return nil, false, fmt.Errorf("'%s': %w", f, err)
		}
		if utf8.Valid(raw) {
			return map[string]any{"content": string(raw)}, true, nil
		}
		return map[string]any{"binaryContent": base64.StdEncoding.EncodeToString(raw)}, true, nil
	case map[string]any:
		source, ok := f["source"].(string)
		if !ok || source == "" {
			return file, false, nil
		}
		raw, err := readFile(baseDir, source)
		if err != nil {
			return nil, false, fmt.Errorf("'%s': %w", source, err)
		}
		delete(f, "source")
		if utf8.Valid(raw) {
			f["content"] = string(raw)
		} else {
			f["binaryContent"] = base64.StdEncoding.EncodeToString(raw)
		}
		return f, true, nil
	}

	v := reflect.ValueOf(file)
	if !v.IsValid() {
		return file, false, nil
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() || v.Elem().Kind() != reflect.Struct {
			return file, false, nil
		}
		copy := reflect.New(v.Elem().Type())
		copy.Elem().Set(v.Elem())
		changed, err := normalizeContainerFileValue(copy.Elem(), baseDir)
		if err != nil || !changed {
			return file, changed, err
		}
		return copy.Interface(), true, nil
	}

	if v.Kind() != reflect.Struct {
		return file, false, nil
	}

	copy := reflect.New(v.Type()).Elem()
	copy.Set(v)
	changed, err := normalizeContainerFileValue(copy, baseDir)
	if err != nil || !changed {
		return file, changed, err
	}
	return copy.Interface(), true, nil
}

func normalizeContainerFileValue(v reflect.Value, baseDir string) (bool, error) {
	sourceField := v.FieldByName("Source")
	if !sourceField.IsValid() || sourceField.Kind() != reflect.Ptr || sourceField.IsNil() {
		return false, nil
	}

	sourceValue := sourceField.Elem()
	if sourceValue.Kind() != reflect.String {
		return false, nil
	}

	source := sourceValue.String()
	raw, err := readFile(baseDir, source)
	if err != nil {
		return false, fmt.Errorf("'%s': %w", source, err)
	}

	sourceField.Set(reflect.Zero(sourceField.Type()))
	if utf8.Valid(raw) {
		return true, setStringPtrField(v.FieldByName("Content"), string(raw))
	}
	return true, setStringPtrField(v.FieldByName("BinaryContent"), base64.StdEncoding.EncodeToString(raw))
}

func setStringPtrField(field reflect.Value, value string) error {
	if !field.IsValid() || field.Kind() != reflect.Ptr || field.Type().Elem().Kind() != reflect.String {
		return nil
	}
	v := reflect.New(field.Type().Elem())
	v.Elem().SetString(value)
	field.Set(v)
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
