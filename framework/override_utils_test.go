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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDotPathParts(t *testing.T) {
	for _, tc := range []struct {
		Input    string
		Expected []string
	}{
		{"", []string{""}},
		{"a", []string{"a"}},
		{"a.b", []string{"a", "b"}},
		{"a.-1", []string{"a", "-1"}},
		{"a.b\\.c", []string{"a", "b.c"}},
		{"a.b\\\\.c", []string{"a", "b\\", "c"}},
	} {
		t.Run(tc.Input, func(t *testing.T) {
			assert.Equal(t, tc.Expected, ParseDotPathParts(tc.Input))
		})
	}
}

func TestWritePathInStruct(t *testing.T) {
	for _, tc := range []struct {
		Name          string
		Spec          string
		Path          []string
		Delete        bool
		Value         interface{}
		Expected      string
		ExpectedError error
	}{
		{
			Name:     "simple object set",
			Spec:     `{"a":{"b":[{}]}}`,
			Path:     []string{"a", "b", "0", "c"},
			Value:    "hello",
			Expected: `{"a":{"b":[{"c":"hello"}]}}`,
		},
		{
			Name:     "simple object delete",
			Spec:     `{"a":{"b":[{"c":"hello"}]}}`,
			Path:     []string{"a", "b", "0", "c"},
			Delete:   true,
			Expected: `{"a":{"b":[{}]}}`,
		},
		{
			Name:     "simple array set",
			Spec:     `{"a":[{}]}`,
			Path:     []string{"a", "0"},
			Value:    "hello",
			Expected: `{"a":["hello"]}`,
		},
		{
			Name:     "simple array append",
			Spec:     `{"a":["hello"]}`,
			Path:     []string{"a", "-1"},
			Value:    "world",
			Expected: `{"a":["hello","world"]}`,
		},
		{
			Name:     "simple array delete",
			Spec:     `{"a":["hello", "world"]}`,
			Path:     []string{"a", "0"},
			Delete:   true,
			Expected: `{"a":["world"]}`,
		},
		{
			Name:     "build object via path",
			Spec:     `{}`,
			Path:     []string{"a", "b"},
			Value:    "hello",
			Expected: `{"a":{"b":"hello"}}`,
		},
		{
			Name:          "bad index str",
			Spec:          `{"a":[]}`,
			Path:          []string{"a", "b"},
			Value:         "hello",
			ExpectedError: fmt.Errorf("a: failed to parse 'b' as array index"),
		},
		{
			Name:          "index out of range",
			Spec:          `{"a": [0]}`,
			Path:          []string{"a", "2"},
			Value:         "hello",
			ExpectedError: fmt.Errorf("a: cannot set '2' in array: out of range"),
		},
		{
			Name:     "no append nested arrays",
			Spec:     `{"a":[[0]]}`,
			Path:     []string{"a", "0", "-1"},
			Value:    "hello",
			Expected: `{"a":[[0, "hello"]]}`,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var inSpec map[string]interface{}
			assert.NoError(t, json.Unmarshal([]byte(tc.Spec), &inSpec))
			outSpec, err := OverridePathInMap(inSpec, tc.Path, tc.Delete, tc.Value)
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error())
				assert.Equal(t, outSpec, map[string]interface{}(nil))
			} else {
				if assert.NoError(t, err) {
					raw, _ := json.Marshal(outSpec)
					assert.JSONEq(t, tc.Expected, string(raw))

					// verify in spec was not modified
					var inSpec2 map[string]interface{}
					assert.NoError(t, json.Unmarshal([]byte(tc.Spec), &inSpec2))
					assert.Equal(t, inSpec, inSpec2)
				}
			}
		})
	}
}

func TestOverrideMapInMap(t *testing.T) {
	input := map[string]interface{}{
		"a": "42",
		"b": []interface{}{"c", "d"},
		"c": map[string]interface{}{
			"d": "42",
			"e": map[string]interface{}{
				"f": "something",
			},
			"g": "other",
		},
		"h": "thing",
	}
	stashInput, _ := json.Marshal(input)
	output, err := OverrideMapInMap(input, map[string]interface{}{
		"a": "13",
		"b": []interface{}{},
		"c": map[string]interface{}{
			"e": map[string]interface{}{
				"z": "thing",
			},
		},
	})
	if assert.NoError(t, err) {
		assert.Equal(t, map[string]interface{}{
			"a": "13",
			"b": []interface{}{},
			"c": map[string]interface{}{
				"d": "42",
				"e": map[string]interface{}{
					"f": "something",
					"z": "thing",
				},
				"g": "other",
			},
			"h": "thing",
		}, output)

		// verify input was not modified
		var inSpec2 map[string]interface{}
		assert.NoError(t, json.Unmarshal(stashInput, &inSpec2))
		assert.Equal(t, input, inSpec2)
	}
}
