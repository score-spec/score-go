package loader

import (
	"errors"
	"io"
	"testing"

	"github.com/score-spec/score-go/types"
	"github.com/stretchr/testify/assert"
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
