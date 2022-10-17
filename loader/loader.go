package loader

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/score-spec/score-go/types"
	"gopkg.in/yaml.v3"
)

// ParseYAML parses YAML into the target mapping structure.
func ParseYAML(src []byte, dest *map[string]interface{}) error {
	return yaml.Unmarshal(src, dest)
}

// MapSpec converts the source mapping structure into the target WorkloadSpec.
func MapSpec(src map[string]interface{}, dest *types.WorkloadSpec) error {
	mapper, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
	})
	if err != nil {
		return fmt.Errorf("initializing decoder: %w", err)
	}

	return mapper.Decode(src)
}
