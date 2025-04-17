// Copyright 2020 Humanitec
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

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"

	"github.com/score-spec/score-go/types"
)

// ValidateJson validates a json structure read from the given reader source. Generally, you must call
// ApplyCommonUpgradeTransforms on the raw structure first unless the input contains zero deprecated concepts.
// For all validation errors, the returned error would be a *jsonschema.ValidationError.
func ValidateJson(r io.Reader) error {
	var obj map[string]interface{}

	var dec = json.NewDecoder(r)
	dec.UseNumber()
	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("decoding source JSON structure: %w", err)
	}

	return Validate(obj)
}

// ValidateYaml validates a yaml structure read from the given reader source. Generally, you must call
// ApplyCommonUpgradeTransforms on the raw structure first unless the input contains zero deprecated concepts.
// For all validation errors returned error would be a *jsonschema.ValidationError.
func ValidateYaml(r io.Reader) error {
	var obj map[string]interface{}

	var dec = yaml.NewDecoder(r)
	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("decoding source YAML structure: %w", err)
	}

	return Validate(obj)
}

// ValidateSpec validates a workload spec structure by serializing it to yaml and calling ValidateYaml.
func ValidateSpec(spec *types.Workload) error {
	intermediate, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal to yaml: %w", err)
	}
	return ValidateYaml(bytes.NewReader(intermediate))
}

// Validate validates the source structure which should be a decoded map. Generally, you must call
// ApplyCommonUpgradeTransforms on the raw structure first unless the input contains zero deprecated concepts.
// For all validation errors returned error would be a *jsonschema.ValidationError.
func Validate(src map[string]interface{}) error {
	schema, err := jsonschema.CompileString("", ScoreSchemaV1b1)
	if err != nil {
		return fmt.Errorf("compiling Score schema: %w", err)
	}

	return schema.Validate(src)
}

// ApplyCommonUpgradeTransforms when we fix aspects of the score spec over time, we sometimes need to break compatibility.
// To reduce affects on users, we can apply a sequence of transformations to the yaml decoded structure so that we can
// fix things on their behalf. This reduces the impact on existing workflows. This function returns messages regarding
// any changes it has made or an error if the structure was unexpected.
// NOTE: this method should only be used for tools or utilities where there is already an established use-case and
// workflow for example score-compose and score-humanitec.
func ApplyCommonUpgradeTransforms(rawScore map[string]interface{}) ([]string, error) {
	changes := make([]string, 0)

	if containersStruct, ok := rawScore["containers"].(map[string]interface{}); ok {
		for name, rawContainerStruct := range containersStruct {
			containerStruct, ok := rawContainerStruct.(map[string]interface{})
			if !ok {
				continue
			}

			// We no longer support multi-line content. Update any arrays in line to be newline-separated
			if filesStruct, ok := containerStruct["files"].([]interface{}); ok {
				filesAsMap := make(map[string]interface{})
				for i, rawFileStruct := range filesStruct {
					fileStruct, ok := rawFileStruct.(map[string]interface{})
					if !ok {
						continue
					}
					if before, ok := fileStruct["content"].([]interface{}); ok {
						delete(fileStruct, "content")
						sb := new(strings.Builder)
						for il, line := range before {
							if il > 0 {
								sb.WriteRune('\n')
							}
							sb.WriteString(fmt.Sprint(line))
						}
						fileStruct["content"] = sb.String()
						changes = append(changes, fmt.Sprintf("containers.%s.files.%d.content: converted from array", name, i))
					}

					target, ok := fileStruct["target"].(string)
					if !ok {
						return nil, fmt.Errorf("containers.%s.files.%d.target: is missing or is not a string", name, i)
					}
					delete(fileStruct, "target")
					filesAsMap[target] = fileStruct
				}

				containerStruct["files"] = filesAsMap
				changes = append(changes, fmt.Sprintf("containers.%s.files: migrated to object", name))
			}

			// We have fixed the naming of the read_only field. It is now readOnly.
			if volumesStruct, ok := containerStruct["volumes"].([]interface{}); ok {
				volumesAsMap := make(map[string]interface{})
				for i, rawVolumeStruct := range volumesStruct {
					volumeStruct, ok := rawVolumeStruct.(map[string]interface{})
					if !ok {
						continue
					}
					if before, ok := volumeStruct["read_only"].(bool); ok {
						delete(volumeStruct, "read_only")
						volumeStruct["readOnly"] = before
						changes = append(changes, fmt.Sprintf("containers.%s.volumes.%d.read_only: migrated to readOnly", name, i))
					}

					target, ok := volumeStruct["target"].(string)
					if !ok {
						return nil, fmt.Errorf("containers.%s.volumes.%d.target: is missing or is not a string", name, i)
					}
					delete(volumeStruct, "target")
					volumesAsMap[target] = volumeStruct
				}

				containerStruct["volumes"] = volumesAsMap
				changes = append(changes, fmt.Sprintf("containers.%s.volumes: migrated to object", name))
			}
		}
	}

	return changes, nil
}
