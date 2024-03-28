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

package framework

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// ParseDotPathParts will parse a common .-separated override path into path elements to traverse.
func ParseDotPathParts(input string) []string {
	// support escaping dot's to insert elements with a . in them.
	input = strings.ReplaceAll(input, "\\\\", "\x01")
	input = strings.ReplaceAll(input, "\\.", "\x00")
	parts := strings.Split(input, ".")
	for i, part := range parts {
		part = strings.ReplaceAll(part, "\x00", ".")
		part = strings.ReplaceAll(part, "\x01", "\\")
		parts[i] = part
	}
	return parts
}

// OverrideMapInMap will take in a decoded json or yaml struct and merge an override map into it. Any maps are merged
// together, other value types are replaced. Nil values will delete overridden keys or otherwise are ignored. This
// returns a shallow copy of the map in a copy-on-write way, so only modified elements are copied.
func OverrideMapInMap(input map[string]interface{}, overrides map[string]interface{}) (map[string]interface{}, error) {
	output := maps.Clone(input)
	for key, value := range overrides {
		if value == nil {
			delete(output, key)
			continue
		}

		existing, hasExisting := output[key]
		if !hasExisting {
			output[key] = value
			continue
		}

		eMap, isEMap := existing.(map[string]interface{})
		vMap, isVMap := value.(map[string]interface{})
		if isEMap && isVMap {
			output[key], _ = OverrideMapInMap(eMap, vMap)
		} else {
			output[key] = value
		}
	}
	return output, nil
}

// OverridePathInMap will take in a decoded json or yaml struct and override a particular path within it with either
// a new value or deletes it. This returns a shallow copy of the map in a copy-on-write way, so only modified elements
// are copied.
func OverridePathInMap(input map[string]interface{}, path []string, isDelete bool, value interface{}) (map[string]interface{}, error) {
	return overridePathInMap(input, path, isDelete, value)
}

func overridePathInMap(input map[string]interface{}, path []string, isDelete bool, value interface{}) (map[string]interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("cannot change root node")
	}

	output := maps.Clone(input)
	if len(path) == 1 {
		if isDelete || value == nil {
			delete(output, path[0])
		} else {
			output[path[0]] = value
		}
		return output, nil
	}

	if _, ok := output[path[0]]; !ok {
		next := make(map[string]interface{})
		subOutput, err := overridePathInMap(next, path[1:], isDelete, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path[0], err)
		}
		output[path[0]] = subOutput
		return output, nil
	}

	switch typed := output[path[0]].(type) {
	case map[string]interface{}:
		subOutput, err := overridePathInMap(typed, path[1:], isDelete, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path[0], err)
		}
		output[path[0]] = subOutput
		return output, nil
	case []interface{}:
		subOutput, err := overridePathInArray(typed, path[1:], isDelete, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path[0], err)
		}
		output[path[0]] = subOutput
		return output, nil
	default:
		return nil, fmt.Errorf("%s: cannot set path in non-map/non-array", path[0])
	}
}

func overridePathInArray(input []interface{}, path []string, isDelete bool, value interface{}) ([]interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("cannot change root node")
	}

	pathIndex, err := strconv.Atoi(path[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse '%s' as array index", path[0])
	}

	output := slices.Clone(input)
	if len(path) == 1 {
		if isDelete || value == nil {
			if pathIndex < 0 || pathIndex >= len(input) {
				return nil, fmt.Errorf("cannot remove '%d' in array: out of range", pathIndex)
			}
			return slices.Delete(output, pathIndex, pathIndex+1), nil
		}
		if pathIndex == -1 {
			output = append(output, value)
			return output, nil
		}
		if pathIndex < 0 || pathIndex >= len(input) {
			return nil, fmt.Errorf("cannot set '%d' in array: out of range", pathIndex)
		}
		output[pathIndex] = value
		return output, nil
	}

	switch typed := output[pathIndex].(type) {
	case map[string]interface{}:
		subOutput, err := overridePathInMap(typed, path[1:], isDelete, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path[0], err)
		}
		output[pathIndex] = subOutput
		return output, nil
	case []interface{}:
		subOutput, err := overridePathInArray(typed, path[1:], isDelete, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path[0], err)
		}
		output[pathIndex] = subOutput
		return output, nil
	default:
		return nil, fmt.Errorf("%s: cannot set path in non-map/non-array", path[0])
	}
}
