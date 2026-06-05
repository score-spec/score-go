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
	"errors"
	"io"
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
			Name: "Embeds source files",
			Input: &types.Workload{
				ApiVersion: "score.dev/v1b1",
				Metadata: types.WorkloadMetadata{
					"name": "hello-world",
				},
				Containers: types.WorkloadContainers{
					"hello": types.Container{
						Files: types.ContainerFiles{
							"/etc/hello-world/config.yaml": types.ContainerFile{
								Source:   stringRef("./test_file.txt"),
								Mode:     stringRef("666"),
								NoExpand: boolRef(true),
							},
							"/etc/hello-world/binary": types.ContainerFile{
								Source: stringRef("./test_binary_file"),
							},
							"/etc/hello-world/short.txt": stringRef("short content"),
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
						Files: types.ContainerFiles{
							"/etc/hello-world/config.yaml": types.ContainerFile{
								Mode:     stringRef("666"),
								Content:  stringRef("Hello World\n"),
								NoExpand: boolRef(true),
							},
							"/etc/hello-world/binary": types.ContainerFile{
								BinaryContent: stringRef("XVLOjEyq5FKgHDGMAYMdp+crq4I="),
							},
							"/etc/hello-world/short.txt": stringRef("short content"),
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
						Files: types.ContainerFiles{
							"/etc/hello-world/config.yaml": types.ContainerFile{
								Source:   stringRef("./not_existing.txt"),
								Mode:     stringRef("666"),
								NoExpand: boolRef(true),
							},
						},
					},
				},
			},
			Error: errors.New("embedding file '/etc/hello-world/config.yaml' for container 'hello': './not_existing.txt': open fixtures/not_existing.txt: no such file or directory"),
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
