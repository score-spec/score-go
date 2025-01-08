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

package loader

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/score-spec/score-go/types"
)

func TestNormalize(t *testing.T) {
	var tests = []struct {
		Name   string
		Source io.Reader
		Input  *types.Workload
		Output *types.Workload
		Error  error
	}{
		{
			Name: "Embeds source file",
			Input: &types.Workload{
				ApiVersion: "score.dev/v1b1",
				Metadata: types.WorkloadMetadata{
					"name": "hello-world",
				},
				Containers: types.WorkloadContainers{
					"hello": types.Container{
						Files: []types.ContainerFilesElem{
							{
								Source:   stringRef("./test_file.txt"),
								Target:   "/etc/hello-world/config.yaml",
								Mode:     stringRef("666"),
								NoExpand: boolRef(true),
							},
						},
					},
				},
			},
			Output: &types.Workload{
				ApiVersion: "score.dev/v1b1",
				Metadata: types.WorkloadMetadata{
					"name": "hello-world",
				},
				Containers: types.WorkloadContainers{
					"hello": types.Container{
						Files: []types.ContainerFilesElem{
							{
								Target:   "/etc/hello-world/config.yaml",
								Mode:     stringRef("666"),
								Content:  stringRef("Hello World\n"),
								NoExpand: boolRef(true),
							},
						},
					},
				},
			},
		},
		{
			Name: "Errors when the source file does not exist",
			Input: &types.Workload{
				ApiVersion: "score.dev/v1b1",
				Metadata: types.WorkloadMetadata{
					"name": "hello-world",
				},
				Containers: types.WorkloadContainers{
					"hello": types.Container{
						Files: []types.ContainerFilesElem{
							{
								Source:   stringRef("./not_existing.txt"),
								Target:   "/etc/hello-world/config.yaml",
								Mode:     stringRef("666"),
								NoExpand: boolRef(true),
							},
						},
					},
				},
			},
			Error: errors.New("embedding file './not_existing.txt' for container 'hello': open fixtures/not_existing.txt: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var err = Normalize(tt.Input, "./fixtures")
			if tt.Error != nil {
				assert.EqualError(t, err, tt.Error.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Output, tt.Input)
			}
		})
	}
}

func TestNormalizeBinaryFile(t *testing.T) {
	td := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(td, "binary"), []byte{0x30, 0x82, 0x03, 0x6b, 0x30}, 0644))
	wkld := &types.Workload{
		ApiVersion: "score.dev/v1b1",
		Metadata: types.WorkloadMetadata{
			"name": "hello-world",
		},
		Containers: types.WorkloadContainers{
			"hello": types.Container{
				Files: []types.ContainerFilesElem{
					{
						Source: stringRef("./binary"),
						Target: "/binary",
					},
				},
			},
		},
	}

	assert.NoError(t, Normalize(wkld, td))
	assert.Equal(t, "0\u0082\x03k0", *wkld.Containers["hello"].Files[0].Content)
	raw, _ := json.Marshal(wkld)
	assert.Equal(t, "{\"apiVersion\":\"score.dev/v1b1\",\"containers\":{\"hello\":{\"files\":[{\"content\":\"0\u0082\\u0003k0\",\"target\":\"/binary\"}],\"image\":\"\"}},\"metadata\":{\"name\":\"hello-world\"}}", string(raw))

	var s map[string]interface{}
	assert.NoError(t, json.Unmarshal(raw, &s))
	assert.Equal(t, "0\u0082\x03k0", s["containers"].(map[string]interface{})["hello"].(map[string]interface{})["files"].([]interface{})[0].(map[string]interface{})["content"])
}
