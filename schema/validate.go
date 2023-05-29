/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package schema

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

var ErrInvalid = errors.New("invalid document")

func ValidateJson(data []byte) error {
	src := gojsonschema.NewStringLoader(string(data))
	return validate(src)
}

func ValidateYaml(data []byte) error {
	var obj map[string]interface{}
	var yamlReader = bytes.NewReader(data)
	if err := yaml.NewDecoder(yamlReader).Decode(&obj); err != nil {
		return fmt.Errorf("decoding yaml: %w", err)
	}

	src := gojsonschema.NewGoLoader(obj)
	return validate(src)
}

func validate(src gojsonschema.JSONLoader) error {
	schema := gojsonschema.NewStringLoader(ScoreSchemaV1b1)

	result, err := gojsonschema.Validate(schema, src)
	if result == nil {
		return fmt.Errorf("validating schema: %w", err)
	}
	if result != nil && !result.Valid() {
		var messages = make([]string, 0)
		var errors = result.Errors()
		for _, err := range errors {
			messages = append(messages, err.String())
		}
		return fmt.Errorf("%v: %s", ErrInvalid, strings.Join(messages, ": "))
	}
	return nil
}
