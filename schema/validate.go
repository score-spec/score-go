/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

// Validates source JSON file.
//
// For all vaidation errors returned error would be a *jsonschema.ValidationError.
func ValidateJson(r io.Reader) error {
	var obj interface{}

	var dec = json.NewDecoder(r)
	dec.UseNumber()
	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("decoding source JSON structure: %w", err)
	}

	return Validate(obj)
}

// Validates source YAML file.
//
// For all vaidation errors returned error would be a *jsonschema.ValidationError.
func ValidateYaml(r io.Reader) error {
	var obj interface{}

	var dec = yaml.NewDecoder(r)
	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("decoding source YAML structure: %w", err)
	}

	return Validate(obj)
}

// Validates source structure.
//
// For all vaidation errors returned error would be a *jsonschema.ValidationError.
func Validate(src interface{}) error {
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
				}
			}

			// We have fixed the naming of the read_only field. It is now readOnly.
			if volumesStruct, ok := containerStruct["volumes"].([]interface{}); ok {
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
				}
			}
		}
	}

	return changes, nil
}
