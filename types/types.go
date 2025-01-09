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

package types

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.15.0 -v --schema-output=https://score.dev/schemas/score=types.gen.go --schema-package=https://score.dev/schemas/score=types --schema-root-type=https://score.dev/schemas/score=Workload ../schema/files/score-v1b1.json.modified

func (m *ResourceMetadata) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out map[string]interface{}
	if err := unmarshal(&out); err != nil {
		return err
	}
	*m = ResourceMetadata(out)
	return nil
}

func (m *WorkloadMetadata) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out map[string]interface{}
	if err := unmarshal(&out); err != nil {
		return err
	}
	*m = WorkloadMetadata(out)
	return nil
}

// findIPowerSuffix is used in ParseResourceLimits to parse memory units.
func findIPowerSuffix(raw string, suffi []string, base int64) (suffix string, m int64) {
	m = 1
	for _, s := range suffi {
		m *= base
		if strings.HasSuffix(raw, s) {
			return s, m
		}
	}
	return "", 0
}

// ParseResourceLimits parses a resource limits definition into milli-cpus and memory bytes if present.
// For example, 500m cpus = 500 millicpus, while 2 cpus = 2000 cpus. 1M == 1000000 bytes of memory, while 1Ki = 1024.
func ParseResourceLimits(rl ResourcesLimits) (milliCpus *int, memoryBytes *int64, err error) {
	if rl.Cpu != nil {
		isMilli := strings.HasSuffix(*rl.Cpu, "m")
		c := strings.TrimSuffix(*rl.Cpu, "m")
		v, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse cpus '%s' as a number", c)
		}
		if !isMilli {
			v *= 1000
		}
		iv := int(v)
		milliCpus = &iv
	}
	if rl.Memory != nil {
		// https://kubernetes.io/docs/tasks/configure-pod-container/assign-memory-resource/#memory-units
		raw := *rl.Memory
		var multiplier int64 = 1
		if s, m := findIPowerSuffix(raw, []string{"K", "M", "G", "T"}, 1000); m > 0 {
			raw = strings.TrimSuffix(raw, s)
			multiplier = m
		} else if s, m = findIPowerSuffix(raw, []string{"Ki", "Mi", "Gi", "Ti"}, 1024); m > 0 {
			raw = strings.TrimSuffix(raw, s)
			multiplier = m
		}
		if v, err := strconv.ParseInt(raw, 10, 64); err != nil {
			return nil, nil, fmt.Errorf("failed to parse memory '%s' as a number", raw)
		} else {
			v *= multiplier
			memoryBytes = &v
		}
	}
	return
}
