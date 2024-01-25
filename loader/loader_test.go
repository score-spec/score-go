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

	"github.com/stretchr/testify/assert"

	"github.com/score-spec/score-go/v1/types"
)

func stringRef(input string) *string {
	return &input
}

func intRef(input int) *int {
	return &input
}

func boolRef(input bool) *bool {
	return &input
}

func schemeRef(input types.HttpProbeScheme) *types.HttpProbeScheme {
	return &input
}

func TestDecodeYaml(t *testing.T) {
	var tests = []struct {
		Name   string
		Source io.Reader
		Output types.Workload
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
apiVersion: score.dev/v1b1
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
      content: |
        ---
        ${resources.env.APP_CONFIG}
      noExpand: true
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
        host: "1.1.1.1"
        scheme: HTTPS
        path: /ready
        port: 8080
        httpHeaders:
        - name: Custom-Header
          value: Awesome

resources:
  env:
    type: environment
  dns:
    type: dns
    class: sensitive
  data:
    type: volume
  db:
    type: postgres
    class: large
    metadata:
      annotations:
        "my.org/version": "0.1"
    params: {
      extensions: {
        uuid-ossp: {
          schema: "uuid_schema",
          version: "1.1"
        }
      }
    }
`)),
			Output: types.Workload{
				ApiVersion: "score.dev/v1b1",
				Metadata: types.WorkloadMetadata{
					Name: "hello-world",
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
						Image:   "busybox",
						Command: []string{"/bin/echo"},
						Args:    []string{"Hello $(FRIEND)"},
						Variables: map[string]string{
							"FRIEND": "World!",
						},
						Files: []types.ContainerFilesElem{
							{
								Target:   "/etc/hello-world/config.yaml",
								Mode:     stringRef("666"),
								Content:  "---\n${resources.env.APP_CONFIG}\n",
								NoExpand: boolRef(true),
							},
						},
						Volumes: []types.ContainerVolumesElem{
							{
								Source:   "${resources.data}",
								Path:     stringRef("sub/path"),
								Target:   "/mnt/data",
								ReadOnly: boolRef(true),
							},
						},
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
								Port: intRef(8080),
							},
						},
						ReadinessProbe: &types.ContainerProbe{
							HttpGet: &types.HttpProbe{
								Host:   stringRef("1.1.1.1"),
								Scheme: schemeRef(types.HttpProbeSchemeHTTPS),
								Path:   "/ready",
								Port:   intRef(8080),
								HttpHeaders: []types.HttpProbeHttpHeadersElem{
									{Name: stringRef("Custom-Header"), Value: stringRef("Awesome")},
								},
							},
						},
					},
				},
				Resources: types.WorkloadResources{
					"env": {
						Type: "environment",
					},
					"dns":  {Type: "dns", Class: stringRef("sensitive")},
					"data": {Type: "volume"},
					"db": {
						Type:  "postgres",
						Class: stringRef("large"),
						Metadata: &types.ResourceMetadata{
							Annotations: map[string]string{
								"my.org/version": "0.1",
							},
						},
						Params: map[string]interface{}{
							"extensions": map[string]interface{}{
								"uuid-ossp": map[string]interface{}{
									"schema":  "uuid_schema",
									"version": "1.1",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var srcMap map[string]interface{}
			var spec types.Workload

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
