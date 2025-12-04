package loader

import (
	"testing"

	"github.com/score-spec/score-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func workloadWith(files types.ContainerFiles, variables types.ContainerVariables, volumes types.ContainerVolumes, resources types.WorkloadResources) *types.Workload {
	return &types.Workload{
		ApiVersion: "score.dev/v1b1",
		Metadata: types.WorkloadMetadata{
			"name": "hello-world",
		},
		Service: &types.WorkloadService{
			Ports: types.WorkloadServicePorts{
				"www": types.ServicePort{
					Port:       80,
					TargetPort: intRef(8080),
				},
			},
		},
		Containers: types.WorkloadContainers{
			"hello": types.Container{
				Image:     "busybox",
				Command:   []string{"/bin/echo"},
				Args:      []string{"Hello $(FRIEND)"},
				Variables: variables,
				Files:     files,
				Volumes:   volumes,
				Resources: &types.ContainerResources{
					Limits: &types.ResourcesLimits{
						Memory: stringRef("128Mi"),
						Cpu:    stringRef("500m"),
					},
					Requests: &types.ResourcesLimits{
						Memory: stringRef("64Mi"),
						Cpu:    stringRef("250m"),
					},
				},
				LivenessProbe: &types.ContainerProbe{
					HttpGet: &types.HttpProbe{
						Path: "/alive",
						Port: 8080,
					},
					Exec: &types.ExecProbe{
						Command: []string{"echo", "hello"},
					},
				},
				ReadinessProbe: &types.ContainerProbe{
					HttpGet: &types.HttpProbe{
						Host:   stringRef("1.1.1.1"),
						Scheme: schemeRef(types.HttpProbeSchemeHTTPS),
						Path:   "/ready",
						Port:   8080,
						HttpHeaders: []types.HttpProbeHttpHeadersElem{
							{Name: "Custom-Header", Value: "Awesome"},
						},
					},
				},
			},
		},
		Resources: resources,
	}
}

func TestValidatePlaceholders(t *testing.T) {
	testCases := []struct {
		name          string
		files         types.ContainerFiles
		variables     types.ContainerVariables
		volumes       types.ContainerVolumes
		resources     types.WorkloadResources
		errorContains []string
	}{
		{
			name: "valid",
			files: types.ContainerFiles{
				"/usr/local/one": {
					Content: stringRef("Placeholder ${resources.res-one.value}"),
				},
				"/usr/local/two": {
					Content: stringRef("Placeholder ${resources.res-two.value} ${resources.res-two.other}"),
				},
				"/usr/local/three": {
					Content: stringRef("No placeholders"),
				},
				"/usr/local/four": {
					Content: stringRef("Escaped $${plaholder}"),
				},
				"/usr/local/five": {
					Content:  stringRef("Invalid placeholder with NoExpand: ${this is invalid}"),
					NoExpand: boolRef(true),
				},
			},
			variables: types.ContainerVariables{
				"VAR_ONE":   "Placeholder ${resources.res-one.value}",
				"VAR_TWO":   "Placeholder ${resources.res-two.value} ${resources.res-two.other}",
				"VAR_THREE": "No placeholders",
				"VAR_FOUR":  "Escaped $${resources.no-exists.value}",
			},
			volumes: types.ContainerVolumes{
				"/mnt/one": {
					Source: "${resources.res-one}",
				},
			},
			resources: types.WorkloadResources{
				"res-one": {
					Type: "type-one",
				},
				"res-two": {
					Type: "type-two",
					Params: types.ResourceParams{
						"var": "${resources.res-one.value}",
					},
				},
			},
		},
		{
			name: "invalid placeholder",
			variables: types.ContainerVariables{
				"INVALID": "Placeholder ${resources.res-one.this has spaces!}",
			},
			errorContains: []string{
				"${resources.res-one.this has spaces!} is malformed",
			},
		},
		{
			name: "resource placeholder with no resource",
			variables: types.ContainerVariables{
				"INVALID": "Placeholder ${resources.res-one.value}",
			},
			errorContains: []string{
				"${resources.res-one.value} does not resolve to a resource",
			},
		},
		{
			name: "invalid first element",
			variables: types.ContainerVariables{
				"INVALID": "Placeholder ${cheese.res-one.value}",
			},
			errorContains: []string{
				"${cheese.res-one.value} has unknown first element",
			},
		},
		{
			name: "one element",
			variables: types.ContainerVariables{
				"INVALID": "Placeholder ${resources}",
			},
			errorContains: []string{
				"${resources} is malformed",
			},
		},
		{
			name: "multiple errors",
			files: types.ContainerFiles{
				"/usr/local/one": {
					Content: stringRef("Placeholder ${resources.res-one.value}"),
				},
				"/usr/local/two": {
					Content: stringRef("Placeholder ${resources.no-exist.value}"),
				},
			},
			variables: types.ContainerVariables{
				"VAR_ONE": "Placeholder ${invalid!}",
			},
			volumes: types.ContainerVolumes{
				"/mnt/one": {
					Source: "${resources.another-no-exist}",
				},
			},
			resources: types.WorkloadResources{
				"res-one": {
					Type: "type-one",
				},
				"res-two": {
					Type: "type-two",
					Params: types.ResourceParams{
						"var": "${resources.yet-another-no-exist.value}",
					},
				},
			},
			errorContains: []string{
				"${resources.no-exist.value} does not resolve to a resource",
				"${invalid!} is malformed",
				"${resources.another-no-exist} does not resolve to a resource",
				"${resources.yet-another-no-exist.value} does not resolve to a resource",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			workload := workloadWith(testCase.files, testCase.variables, testCase.volumes, testCase.resources)
			err := Validate(workload)
			if len(testCase.errorContains) == 0 {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			for _, msg := range testCase.errorContains {
				assert.ErrorContains(t, err, msg)
			}
		})
	}
}
