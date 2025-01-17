// Copyright 2020 Humanitec
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

package schema

import (
	"bytes"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func newTestDocument() map[string]interface{} {
	var data = []byte(`
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
      content: "${resources.env.APP_CONFIG}"
    - target: /etc/hello-world/binary
      mode: "755"
      binaryContent: "aGVsbG8="
    volumes:
    - source: ${resources.data}
      path: sub/path
      target: /mnt/data
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
`)

	var obj map[string]interface{}
	var yamlReader = bytes.NewReader(data)
	_ = yaml.NewDecoder(yamlReader).Decode(&obj)
	return obj
}

func TestSchema(t *testing.T) {
	var tests = []struct {
		Name    string
		Src     map[string]interface{}
		Message string
	}{
		{
			Name: "Valid input",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				return src
			}(),
			Message: "",
		},

		// apiVersion
		//
		{
			Name: "apiVersion is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				delete(src, "apiVersion")
				return src
			}(),
			Message: "missing properties: 'apiVersion'",
		},
		{
			Name: "apiVersion is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["apiVersion"] = nil
				return src
			}(),
			Message: "/apiVersion",
		},
		{
			Name: "apiVersion is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["apiVersion"] = 12
				return src
			}(),
			Message: "/apiVersion",
		},

		// metadata
		//
		{
			Name: "metadata is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				delete(src, "metadata")
				return src
			}(),
			Message: "missing properties: 'metadata'",
		},
		{
			Name: "metadata is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["metadata"] = nil
				return src
			}(),
			Message: "/metadata",
		},
		{
			Name: "metadata.name is required",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				delete(src["metadata"].(map[string]interface{}), "name")
				return src
			}(),
			Message: "/metadata",
		},
		{
			Name: "metadata.name is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["metadata"].(map[string]interface{})["name"] = 12
				return src
			}(),
			Message: "/metadata/name",
		},

		// service
		//
		{
			Name: "service is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = nil
				return src
			}(),
			Message: "/service",
		},
		{
			Name: "service.ports is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": nil,
				}
				return src
			}(),
			Message: "/service/ports",
		},
		{
			Name: "service.ports is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{},
				}
				return src
			}(),
			Message: "",
		},
		{
			Name: "service.ports.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": nil,
					},
				}
				return src
			}(),
			Message: "/service/ports/www",
		},
		{
			Name: "service.ports.*.port is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{},
					},
				}
				return src
			}(),
			Message: "/service/ports/www",
		},
		{
			Name: "service.ports.*.port is not a number",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"port":       false,
							"targetPort": 8080,
						},
					},
				}
				return src
			}(),
			Message: "/service/ports/www/port",
		},
		{
			Name: "service.ports.*.targetPort is not a number",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"port":       80,
							"targetPort": false,
						},
					},
				}
				return src
			}(),
			Message: "/service/ports/www/targetPort",
		},
		{
			Name: "service.ports.*.targetPort is optional",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"port": 8080,
						},
					},
				}
				return src
			}(),
			Message: "",
		},
		{
			Name: "service with multiple ports",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"port":       80,
							"targetPort": 8080,
						},
						"admin": map[string]interface{}{
							"port": 90,
						},
					},
				}
				return src
			}(),
			Message: "",
		},

		// containers
		//
		{
			Name: "containers is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				delete(src, "containers")
				return src
			}(),
			Message: "missing properties: 'containers'",
		},
		{
			Name: "containers is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"] = nil
				return src
			}(),
			Message: "/containers",
		},
		{
			Name: "containers is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"] = map[string]interface{}{}
				return src
			}(),
			Message: "/containers",
		},
		{
			Name: "containers.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"].(map[string]interface{})["hello"] = nil
				return src
			}(),
			Message: "/containers/hello",
		},

		// containers.*.image
		//
		{
			Name: "containers.*.image is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				delete(hello, "image")
				return src
			}(),
			Message: "/containers/hello",
		},
		{
			Name: "containers.*.image is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["image"] = nil
				return src
			}(),
			Message: "/containers/hello/image",
		},
		{
			Name: "containers.*.image is not a sring",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["image"] = 12
				return src
			}(),
			Message: "/containers/hello/image",
		},

		// containers.*.command
		//
		{
			Name: "containers.*.command is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["command"] = nil
				return src
			}(),
			Message: "/containers/hello/command",
		},
		{
			Name: "containers.*.command is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["command"] = []interface{}{}
				return src
			}(),
			Message: "",
		},

		// containers.*.args
		//
		{
			Name: "containers.*.args is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["args"] = nil
				return src
			}(),
			Message: "/containers/hello/args",
		},
		{
			Name: "containers.*.args is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["args"] = []interface{}{}
				return src
			}(),
			Message: "",
		},

		// containers.*.variables
		//
		{
			Name: "containers.*.variables is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"] = nil
				return src
			}(),
			Message: "/containers/hello/variables",
		},
		{
			Name: "containers.*.variables is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.variables.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"].(map[string]interface{})["FRIEND"] = nil
				return src
			}(),
			Message: "/containers/hello/variables/FRIEND",
		},
		{
			Name: "containers.*.variables.* is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"].(map[string]interface{})["FRIEND"] = 12
				return src
			}(),
			Message: "/containers/hello/variables/FRIEND",
		},

		// containers.*.files
		//
		{
			Name: "containers.*.files is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["files"] = nil
				return src
			}(),
			Message: "/containers/hello/files",
		},
		{
			Name: "containers.*.files is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["files"] = []interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.files.*.target is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "target")
				return src
			}(),
			Message: "/containers/hello/files/0",
		},
		{
			Name: "containers.*.files.*.target is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["target"] = nil
				return src
			}(),
			Message: "/containers/hello/files/0/target",
		},
		{
			Name: "containers.*.files.*.target is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["target"] = 12
				return src
			}(),
			Message: "/containers/hello/files/0/target",
		},
		{
			Name: "containers.*.files.*.mode is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["mode"] = nil
				return src
			}(),
			Message: "/containers/hello/files/0/mode",
		},
		{
			Name: "containers.*.files.*.mode is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["mode"] = 12
				return src
			}(),
			Message: "/containers/hello/files/0/mode",
		},
		{
			Name: "containers.*.files.*.source is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				return src
			}(),
			Message: "/containers/hello/files/0",
		},
		{
			Name: "containers.*.files.*.source is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				file["source"] = nil
				return src
			}(),
			Message: "/containers/hello/files/0",
		},
		{
			Name: "containers.*.files.*.source is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				file["source"] = ""
				return src
			}(),
			Message: "/containers/hello/files/0/source",
		},
		{
			Name: "containers.*.files.*.source is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				file["source"] = 5
				return src
			}(),
			Message: "/containers/hello/files/0/source",
		},
		{
			Name: "containers.*.files.*.binaryContent is bad format",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				file["binaryContent"] = map[string]interface{}{}
				return src
			}(),
			Message: "/containers/hello/files/0/binaryContent",
		},
		{
			Name: "containers.*.files.*.noExpand is set to true",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["noExpand"] = true
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.files.*.noExpand isset to false",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["noExpand"] = false
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.files.*.noExpand is not a boolean",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				file["noExpand"] = 5
				return src
			}(),
			Message: "/containers/hello/files/0/noExpand",
		},
		{
			Name: "containers.*.files.*.content is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				delete(file, "content")
				return src
			}(),
			Message: "/containers/hello/files/0",
		},
		{
			Name: "containers.*.files.*.content is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["content"] = nil
				return src
			}(),
			Message: "/containers/hello/files/0/content",
		},
		{
			Name: "containers.*.files.*.content is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["content"] = 5
				return src
			}(),
			Message: "/containers/hello/files/0/content",
		},
		{
			Name: "containers.*.files.*.content is an empty array",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["content"] = []interface{}{}
				return src
			}(),
			Message: "/containers/hello/files/0/content",
		},
		{
			Name: "containers.*.files.*.content is an array of strings",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["content"] = []interface{}{
					"Line 1",
					"Line 2",
				}
				return src
			}(),
			Message: "/containers/hello/files/0/content",
		},

		// containers.*.volumes
		//
		{
			Name: "containers.*.volumes is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["volumes"] = nil
				return src
			}(),
			Message: "/containers/hello/volumes",
		},
		{
			Name: "containers.*.volumes is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["volumes"] = []interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.volumes.*.source is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				delete(volumes, "source")
				return src
			}(),
			Message: "/containers/hello/volumes/0",
		},
		{
			Name: "containers.*.volumes.*.source is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["source"] = nil
				return src
			}(),
			Message: "/containers/hello/volumes/0/source",
		},
		{
			Name: "containers.*.volumes.*.source is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["source"] = 12
				return src
			}(),
			Message: "/containers/hello/volumes/0/source",
		},
		{
			Name: "containers.*.volumes.*.path is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["path"] = nil
				return src
			}(),
			Message: "/containers/hello/volumes/0/path",
		},
		{
			Name: "containers.*.volumes.*.path is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["path"] = 12
				return src
			}(),
			Message: "/containers/hello/volumes/0/path",
		},
		{
			Name: "containers.*.volumes.*.target is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				delete(volumes, "target")
				return src
			}(),
			Message: "/containers/hello/volumes/0",
		},
		{
			Name: "containers.*.volumes.*.target is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["target"] = nil
				return src
			}(),
			Message: "/containers/hello/volumes/0/target",
		},
		{
			Name: "containers.*.volumes.*.target is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["target"] = 12
				return src
			}(),
			Message: "/containers/hello/volumes/0/target",
		},
		{
			Name: "containers.*.volumes.*.readOnly is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["readOnly"] = nil
				return src
			}(),
			Message: "/containers/hello/volumes/0/readOnly",
		},
		{
			Name: "containers.*.volumes.*.readOnly is not a boolean",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["readOnly"] = 12
				return src
			}(),
			Message: "/containers/hello/volumes/0/readOnly",
		},

		// containers.*.resources
		//
		{
			Name: "containers.*.resources is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"] = nil
				return src
			}(),
			Message: "/containers/hello/resources",
		},
		{
			Name: "containers.*.resources is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.resources.limits is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["limits"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/limits",
		},
		{
			Name: "containers.*.resources.limits is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["limits"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.resources.limits.memory is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var limits = hello["resources"].(map[string]interface{})["limits"].(map[string]interface{})
				limits["memory"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/limits/memory",
		},
		{
			Name: "containers.*.resources.limits.memory is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var limits = hello["resources"].(map[string]interface{})["limits"].(map[string]interface{})
				limits["memory"] = 12
				return src
			}(),
			Message: "/containers/hello/resources/limits/memory",
		},
		{
			Name: "containers.*.resources.limits.memory is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var limits = hello["resources"].(map[string]interface{})["limits"].(map[string]interface{})
				limits["memory"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/limits/memory",
		},
		{
			Name: "containers.*.resources.limits.cpu is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var limits = hello["resources"].(map[string]interface{})["limits"].(map[string]interface{})
				limits["cpu"] = 12
				return src
			}(),
			Message: "/containers/hello/resources/limits/cpu",
		},
		{
			Name: "containers.*.resources.requests is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["requests"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/requests",
		},
		{
			Name: "containers.*.resources.requests is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["requests"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.resources.requests.memory is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var requests = hello["resources"].(map[string]interface{})["requests"].(map[string]interface{})
				requests["memory"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/requests/memory",
		},
		{
			Name: "containers.*.resources.requests.memory is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var requests = hello["resources"].(map[string]interface{})["requests"].(map[string]interface{})
				requests["memory"] = 12
				return src
			}(),
			Message: "/containers/hello/resources/requests/memory",
		},
		{
			Name: "containers.*.resources.requests.memory is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var requests = hello["resources"].(map[string]interface{})["requests"].(map[string]interface{})
				requests["memory"] = nil
				return src
			}(),
			Message: "/containers/hello/resources/requests/memory",
		},
		{
			Name: "containers.*.resources.requests.cpu is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var requests = hello["resources"].(map[string]interface{})["requests"].(map[string]interface{})
				requests["cpu"] = 12
				return src
			}(),
			Message: "/containers/hello/resources/requests/cpu",
		},

		// containers.*.livenessProbe
		//
		{
			Name: "containers.*.livenessProbe is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe",
		},
		{
			Name: "containers.*.livenessProbe.exec is nil",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"].(map[string]interface{})["exec"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe/exec",
		},
		{
			Name: "containers.*.livenessProbe.exec is bad",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"].(map[string]interface{})["exec"] = map[string]interface{}{
					"command": true,
				}
				return src
			}(),
			Message: "/containers/hello/livenessProbe/exec/command",
		},
		{
			Name: "containers.*.livenessProbe.httpGet is nil",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"].(map[string]interface{})["httpGet"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.host is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["host"] = 12
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/host",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.host is 1.1.1.1",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["host"] = "1.1.1.1"
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.scheme is HTTP",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["scheme"] = "HTTP"
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.scheme is HTTPS",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["scheme"] = "HTTP"
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.scheme is TCP",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["scheme"] = "TCP"
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/scheme",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.scheme is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["scheme"] = 12
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/scheme",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.path is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				delete(httpGet, "path")
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.path is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["path"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/path",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.path is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["path"] = 12
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/path",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.port is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["port"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/port",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.port is not a number",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["port"] = "12"
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/port",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = nil
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/httpHeaders",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders.*.name is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  nil,
						"value": "Awesome",
					},
				}
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/httpHeaders/0/name",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders.*.name is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  12,
						"value": "Awesome",
					},
				}
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/httpHeaders/0/name",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders.*.value is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  "Custom-Header",
						"value": nil,
					},
				}
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/httpHeaders/0/value",
		},
		{
			Name: "containers.*.livenessProbe.httpGet.httpHeaders.*.value is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["livenessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  "Custom-Header",
						"value": nil,
					},
				}
				return src
			}(),
			Message: "/containers/hello/livenessProbe/httpGet/httpHeaders/0/value",
		},

		// containers.*.readinessProbe
		//
		{
			Name: "containers.*.readinessProbe is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["readinessProbe"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe",
		},
		{
			Name: "containers.*.readinessProbe.exec is nil",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["readinessProbe"].(map[string]interface{})["exec"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe/exec",
		},
		{
			Name: "containers.*.readinessProbe.httpGet is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["readinessProbe"].(map[string]interface{})["httpGet"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.path is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				delete(httpGet, "path")
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.path is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["path"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/path",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.path is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["path"] = 12
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/path",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.port is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["port"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/port",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.port is not a number",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["port"] = "12"
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/port",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = nil
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/httpHeaders",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders.*.name is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  nil,
						"value": "Awesome",
					},
				}
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/httpHeaders/0/name",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders.*.name is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  12,
						"value": "Awesome",
					},
				}
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/httpHeaders/0/name",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders.*.value is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  "Custom-Header",
						"value": nil,
					},
				}
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/httpHeaders/0/value",
		},
		{
			Name: "containers.*.readinessProbe.httpGet.httpHeaders.*.value is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var httpGet = hello["readinessProbe"].(map[string]interface{})["httpGet"].(map[string]interface{})
				httpGet["httpHeaders"] = []interface{}{
					map[string]interface{}{
						"name":  "Custom-Header",
						"value": nil,
					},
				}
				return src
			}(),
			Message: "/containers/hello/readinessProbe/httpGet/httpHeaders/0/value",
		},

		// resources
		//
		{
			Name: "resources is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["resources"] = nil
				return src
			}(),
			Message: "/resources",
		},
		{
			Name: "resources is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["resources"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "resources.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["resources"].(map[string]interface{})["db"] = nil
				return src
			}(),
			Message: "/resources",
		},

		// resources.*.type
		//
		{
			Name: "resources.*.type is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				delete(db, "type")
				return src
			}(),
			Message: "/resources/db",
		},
		{
			Name: "resources.*.type is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["type"] = nil
				return src
			}(),
			Message: "/resources/db/type",
		},
		{
			Name: "resources.*.type is not a sring",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["type"] = 12
				return src
			}(),
			Message: "/resources/db/type",
		},

		// resource.*.class
		{
			Name: "resources.*.class is not valid",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["class"] = "cl@ss?"
				return src
			}(),
			Message: "/resources/db/class",
		},

		// resources.*.metadata
		//
		{
			Name: "resources.*.metadata is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["metadata"] = nil
				return src
			}(),
			Message: "/resources/db/metadata",
		},
		{
			Name: "resources.*.metadata is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["metadata"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},

		// resources.*.metadata.annotations
		//
		{
			Name: "resources.*.metadata.annotations is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				var metadata = db["metadata"].(map[string]interface{})
				metadata["annotations"] = nil
				return src
			}(),
			Message: "/resources/db/metadata/annotations",
		},
		{
			Name: "resources.*.metadata.annotations is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				var metadata = db["metadata"].(map[string]interface{})
				metadata["annotations"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "resources.*.metadata.annotations.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				var metadata = db["metadata"].(map[string]interface{})
				metadata["annotations"] = map[string]interface{}{
					"one": nil,
					"two": "three",
				}
				return src
			}(),
			Message: "/resources/db/metadata/annotations/one",
		},
		{
			Name: "resources.*.metadata.annotations.* is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				var metadata = db["metadata"].(map[string]interface{})
				metadata["annotations"] = map[string]interface{}{
					"one": 12,
					"two": "three",
				}
				return src
			}(),
			Message: "/resources/db/metadata/annotations/one",
		},

		// resources.*.params
		//
		{
			Name: "resources.*.params is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["params"] = nil
				return src
			}(),
			Message: "/resources/db/params",
		},
		{
			Name: "resources.*.params is not an object",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["params"] = 12
				return src
			}(),
			Message: "/resources/db/params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := Validate(tt.Src)

			if tt.Message == "" {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, &jsonschema.ValidationError{}, err)
				assert.ErrorContains(t, err, tt.Message)
			}
		})
	}
}
