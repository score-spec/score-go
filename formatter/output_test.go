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

package formatter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableOutputFormatter_Display(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		rows    [][]string
		want    string
	}{
		{
			name:    "simple table",
			headers: []string{"Name", "Value"},
			rows:    [][]string{{"n1", "v1"}, {"n2", "v2"}},
			want: `+------+-------+
| NAME | VALUE |
+------+-------+
| n1   | v1    |
+------+-------+
| n2   | v2    |
+------+-------+
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f := &TableOutputFormatter{
				Headers: tt.headers,
				Rows:    tt.rows,
				Out:     buf,
			}
			err := f.Display()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestJSONOutputFormatter_Display(t *testing.T) {
	tests := []struct {
		name string
		data []interface{}
		want string
	}{
		{
			name: "simple object",
			data: []interface{}{map[string]string{"k1": "v1"}, map[string]string{"k2": "v2"}},
			want: "[\n  {\n    \"k1\": \"v1\"\n  },\n  {\n    \"k2\": \"v2\"\n  }\n]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f := &JSONOutputFormatter[interface{}]{
				Data: tt.data,
				Out:  buf,
			}
			err := f.Display()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestYAMLOutputFormatter_Display(t *testing.T) {
	tests := []struct {
		name string
		data []interface{}
		want string
	}{
		{
			name: "simple object",
			data: []interface{}{map[string]string{"k1": "v1"}, map[string]string{"k2": "v2"}},
			want: "- k1: v1\n- k2: v2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f := &YAMLOutputFormatter[interface{}]{
				Data: tt.data,
				Out:  buf,
			}
			err := f.Display()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
