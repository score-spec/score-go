// Copyright 2025 The Score Authors
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
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/score-spec/score-go/types"
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
      /etc/hello-world/config.yaml:
        mode: "666"
        content: |
          ---
          ${resources.env.APP_CONFIG}
        noExpand: true
      /etc/hello-world/binary:
        binaryContent: aGVsbG8=
    volumes:
      /mnt/data:
        source: ${resources.data}
        path: sub/path
        readOnly: true
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
      exec:
        command:
        - echo
        - hello
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
						Image:   "busybox",
						Command: []string{"/bin/echo"},
						Args:    []string{"Hello $(FRIEND)"},
						Variables: map[string]string{
							"FRIEND": "World!",
						},
						Files: map[string]types.ContainerFile{
							"/etc/hello-world/config.yaml": {
								Mode:     stringRef("666"),
								Content:  stringRef("---\n${resources.env.APP_CONFIG}\n"),
								NoExpand: boolRef(true),
							},
							"/etc/hello-world/binary": {
								BinaryContent: stringRef("aGVsbG8="),
							},
						},
						Volumes: map[string]types.ContainerVolume{
							"/mnt/data": {
								Source:   "${resources.data}",
								Path:     stringRef("sub/path"),
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
				Resources: types.WorkloadResources{
					"env": {
						Type: "environment",
					},
					"dns":  {Type: "dns", Class: stringRef("sensitive")},
					"data": {Type: "volume"},
					"db": {
						Type:  "postgres",
						Class: stringRef("large"),
						Metadata: map[string]interface{}{
							"annotations": map[string]interface{}{
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

			var err = yaml.NewDecoder(tt.Source).Decode(&srcMap)
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
