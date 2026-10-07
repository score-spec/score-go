// Copyright 2025 The Score Authors
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
	"fmt"
	"io"

	"github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"

	"github.com/score-spec/score-go/types"
)

// ParseYAML parses YAML into the target mapping structure.
// Deprecated. Please use the yaml/v3 library directly rather than calling this method.
func ParseYAML(dest *map[string]interface{}, r io.Reader) error {
	return yaml.NewDecoder(r).Decode(dest)
}

// MapSpec converts the source mapping structure into the target WorkloadSpec.
func MapSpec(dest *types.Workload, src map[string]interface{}) error {
	src = normalizeShortFileAndVolumeForms(src)
	mapper, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
	})
	if err != nil {
		return fmt.Errorf("initializing decoder: %w", err)
	}
	return mapper.Decode(src)
}

func normalizeShortFileAndVolumeForms(src map[string]interface{}) map[string]interface{} {
	containers, ok := src["containers"].(map[string]interface{})
	if !ok {
		return src
	}

	normalizedContainers := make(map[string]interface{}, len(containers))
	containersChanged := false
	for name, rawContainer := range containers {
		container, ok := rawContainer.(map[string]interface{})
		if !ok {
			normalizedContainers[name] = rawContainer
			continue
		}

		var normalizedContainer map[string]interface{}
		for _, field := range []string{"files", "volumes"} {
			shortFormProperty := "source"
			if field == "files" {
				shortFormProperty = "content"
			}
			entries, ok := container[field].(map[string]interface{})
			if !ok {
				continue
			}

			normalizedEntries := make(map[string]interface{}, len(entries))
			entriesChanged := false
			for target, rawEntry := range entries {
				if source, ok := rawEntry.(string); ok {
					normalizedEntries[target] = map[string]interface{}{shortFormProperty: source}
					entriesChanged = true
				} else {
					normalizedEntries[target] = rawEntry
				}
			}
			if entriesChanged {
				if normalizedContainer == nil {
					normalizedContainer = cloneMap(container)
				}
				normalizedContainer[field] = normalizedEntries
			}
		}
		if normalizedContainer != nil {
			normalizedContainers[name] = normalizedContainer
			containersChanged = true
		} else {
			normalizedContainers[name] = rawContainer
		}
	}
	if !containersChanged {
		return src
	}

	normalizedSource := cloneMap(src)
	normalizedSource["containers"] = normalizedContainers
	return normalizedSource
}

func cloneMap(source map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{}, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}
