/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package loader

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/score-spec/score-go/types"
	"github.com/stretchr/testify/assert"
)

func TestDecodeYaml(t *testing.T) {
	var tests = []struct {
		Name   string
		Source io.Reader
		Output types.WorkloadSpec
		Error  error
	}{
		{
			Name:   "Should handle empty input",
			Source: bytes.NewReader([]byte{}),
			Error:  errors.New("EOF"),
		},
		{
			Name:   "Should handle invalid YAML input",
			Source: bytes.NewReader([]byte("<NOT A VALID YAML>")),
			Error:  errors.New("cannot unmarshal"),
		},
		{
			Name: "Should decode the SCORE spec",
			Source: bytes.NewReader([]byte(`
---
apiVersion: score.sh/v1b1
metadata:
  name: hello-world

service:
  ports:
    www:
      port: 80
      targetPort: 8080

containers:
  hello:
    image: busybox
    command: 
    - "/bin/echo"
    args:
    - "Hello $(FRIEND)"
    variables:
      FRIEND: World!
    files:
    - target: /etc/hello-world/config.yaml
      mode: "666"
      content:
      - "---"
      - ${resources.env.APP_CONFIG}
    volumes:
    - source: ${resources.data}
      path: sub/path
      target: /mnt/data
      read_only: true
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
      requests:
        memory: "64Mi"
        cpu: "250m"
    livenessProbe:
      httpGet:
        path: /alive
        port: 8080
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
        httpHeaders:
        - name: Custom-Header
          value: Awesome

resources:
  env:
    type: environment
    properties:
      APP_CONFIG:
  dns:
    type: dns
  data:
    type: volume
  db:
    type: postgres
    properties:
      host:
        type: string
        default: localhost
        required: true
      port:
        default: 5432
      user.name:
`)),
			Output: types.WorkloadSpec{
				ApiVersion: "score.sh/v1b1",
				Metadata: types.WorkloadMeta{
					Name: "hello-world",
				},
				Service: types.ServiceSpec{
					Ports: types.ServicePortsSpecs{
						"www": types.ServicePortSpec{
							Port:       80,
							Protocol:   "",
							TargetPort: 8080,
						},
					},
				},
				Containers: types.ContainersSpecs{
					"hello": types.ContainerSpec{
						Image:   "busybox",
						Command: []string{"/bin/echo"},
						Args:    []string{"Hello $(FRIEND)"},
						Variables: map[string]string{
							"FRIEND": "World!",
						},
						Files: []types.FileMountSpec{
							{
								Target: "/etc/hello-world/config.yaml",
								Mode:   "666",
								Content: []string{
									"---",
									"${resources.env.APP_CONFIG}",
								},
							},
						},
						Volumes: []types.VolumeMountSpec{
							{
								Source:   "${resources.data}",
								Path:     "sub/path",
								Target:   "/mnt/data",
								ReadOnly: true,
							},
						},
						Resources: types.ContainerResourcesRequirementsSpec{
							Limits: map[string]interface{}{
								"memory": "128Mi",
								"cpu":    "500m",
							},
							Requests: map[string]interface{}{
								"memory": "64Mi",
								"cpu":    "250m",
							},
						},
						LivenessProbe: types.ContainerProbeSpec{
							HTTPGet: types.HTTPGetActionSpec{
								Path: "/alive",
								Port: 8080,
							},
						},
						ReadinessProbe: types.ContainerProbeSpec{
							HTTPGet: types.HTTPGetActionSpec{
								Path: "/ready",
								Port: 8080,
								HTTPHeaders: []types.HTTPHeaderSpec{
									{Name: "Custom-Header", Value: "Awesome"},
								},
							},
						},
					},
				},
				Resources: map[string]types.ResourceSpec{
					"env": {
						Type: "environment",
						Properties: map[string]types.ResourcePropertySpec{
							"APP_CONFIG": {},
						},
					},
					"dns":  {Type: "dns"},
					"data": {Type: "volume"},
					"db": {
						Type: "postgres",
						Properties: map[string]types.ResourcePropertySpec{
							"host":      {Type: "string", Default: "localhost", Required: true},
							"port":      {Default: 5432},
							"user.name": {},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var srcMap map[string]interface{}
			var spec types.WorkloadSpec

			var err = ParseYAML(&srcMap, tt.Source)
			if err == nil {
				err = MapSpec(&spec, srcMap)
			}

			if tt.Error != nil {
				// On Error
				//
				assert.ErrorContains(t, err, tt.Error.Error())
			} else {
				// On Success
				//
				assert.NoError(t, err)
				assert.Equal(t, tt.Output, spec)
			}
		})
	}
}
