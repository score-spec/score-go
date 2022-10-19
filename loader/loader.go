package loader

import (
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/score-spec/score-go/types"
	"gopkg.in/yaml.v3"
)

// ParseYAML parses YAML into the target mapping structure.
func ParseYAML(dest *map[string]interface{}, r io.Reader) error {
	return yaml.NewDecoder(r).Decode(dest)
}

// MapSpec converts the source mapping structure into the target WorkloadSpec.
func MapSpec(dest *types.WorkloadSpec, src map[string]interface{}) error {
	mapper, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
	})
	if err != nil {
		return fmt.Errorf("initializing decoder: %w", err)
	}

	return mapper.Decode(src)
}
