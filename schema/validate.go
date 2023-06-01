/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package schema

import (
	"fmt"
	"io"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// Validates source JSON file.
func ValidateJson(r io.Reader) (*gojsonschema.Result, error) {
	src, rdr := gojsonschema.NewReaderLoader(r)
	if _, err := io.ReadAll(rdr); err != nil {
		return nil, fmt.Errorf("reading JSON: %w", err)
	}
	return Validate(src)
}

// Validates source YAML file.
func ValidateYaml(r io.Reader) (*gojsonschema.Result, error) {
	var obj map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&obj); err != nil {
		return nil, fmt.Errorf("decoding YAML: %w", err)
	}

	src := gojsonschema.NewGoLoader(obj)
	return Validate(src)
}

// Validates source Score structure.
func Validate(src gojsonschema.JSONLoader) (*gojsonschema.Result, error) {
	schema := gojsonschema.NewStringLoader(ScoreSchemaV1b1)
	return gojsonschema.Validate(schema, src)
}
