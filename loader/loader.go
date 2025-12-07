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

	"github.com/mitchellh/mapstructure"
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
	mapper, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
	})
	if err != nil {
		return fmt.Errorf("initializing decoder: %w", err)
	}
	return mapper.Decode(src)
}
