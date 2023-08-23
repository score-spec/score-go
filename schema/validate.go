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
