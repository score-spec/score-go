/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package schema

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
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
  dns:
    type: dns
  data:
    type: volume
  db:
    type: postgres
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
	yaml.NewDecoder(yamlReader).Decode(&obj)
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
			Message: "apiVersion is required",
		},
		{
			Name: "apiVersion is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["apiVersion"] = nil
				return src
			}(),
			Message: "apiVersion: Invalid type",
		},
		{
			Name: "apiVersion is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["apiVersion"] = 12
				return src
			}(),
			Message: "apiVersion: Invalid type",
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
			Message: "metadata is required",
		},
		{
			Name: "metadata is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["metadata"] = nil
				return src
			}(),
			Message: "metadata: Invalid type",
		},
		{
			Name: "metadata.name is required",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				delete(src["metadata"].(map[string]interface{}), "name")
				return src
			}(),
			Message: "metadata: name is required",
		},
		{
			Name: "metadata.name is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["metadata"].(map[string]interface{})["name"] = 12
				return src
			}(),
			Message: "metadata.name: Invalid type",
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
			Message: "service: Invalid type",
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
			Message: "service.ports: Invalid type",
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
			Message: "service.ports: Must have at least 1 properties",
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
			Message: "service.ports.www: Invalid type",
		},
		{
			Name: "service.ports.*.targetPort is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{},
					},
				}
				return src
			}(),
			Message: "service.ports.www: targetPort is required",
		},
		{
			Name: "service.ports.*.targetPort is not a number",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"targetPort": false,
						},
					},
				}
				return src
			}(),
			Message: "service.ports.www.targetPort: Invalid type",
		},
		{
			Name: "service.ports.*.port is optional",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["service"] = map[string]interface{}{
					"ports": map[string]interface{}{
						"www": map[string]interface{}{
							"targetPort": 8080,
						},
					},
				}
				return src
			}(),
			Message: "",
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
			Message: "service.ports.www.port: Invalid type",
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
							"targetPort": 8090,
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
			Message: "containers is required",
		},
		{
			Name: "containers is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"] = nil
				return src
			}(),
			Message: "containers: Invalid type",
		},
		{
			Name: "containers is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers: Must have at least 1 properties",
		},
		{
			Name: "containers.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["containers"].(map[string]interface{})["hello"] = nil
				return src
			}(),
			Message: "containers.hello: Invalid type",
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
			Message: "containers.hello: image is required",
		},
		{
			Name: "containers.*.image is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["image"] = nil
				return src
			}(),
			Message: "containers.hello.image: Invalid type",
		},
		{
			Name: "containers.*.image is not a sring",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["image"] = 12
				return src
			}(),
			Message: "containers.hello.image: Invalid type",
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
			Message: "containers.hello.command: Invalid type",
		},
		{
			Name: "containers.*.command is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["command"] = []string{}
				return src
			}(),
			Message: "containers.hello.command: Array must have at least 1 items",
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
			Message: "containers.hello.args: Invalid type",
		},
		{
			Name: "containers.*.args is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["args"] = []string{}
				return src
			}(),
			Message: "containers.hello.args: Array must have at least 1 items",
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
			Message: "containers.hello.variables: Invalid type",
		},
		{
			Name: "containers.*.variables is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.variables: Must have at least 1 properties",
		},
		{
			Name: "containers.*.variables.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"].(map[string]interface{})["FRIEND"] = nil
				return src
			}(),
			Message: "containers.hello.variables.FRIEND: Invalid type",
		},
		{
			Name: "containers.*.variables.* is not a string",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["variables"].(map[string]interface{})["FRIEND"] = 12
				return src
			}(),
			Message: "containers.hello.variables.FRIEND: Invalid type",
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
			Message: "containers.hello.files: Invalid type",
		},
		{
			Name: "containers.*.files is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["files"] = []interface{}{}
				return src
			}(),
			Message: "containers.hello.files: Array must have at least 1 items",
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
			Message: "containers.hello.files.0: target is required",
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
			Message: "containers.hello.files.0.target: Invalid type",
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
			Message: "containers.hello.files.0.target: Invalid type",
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
			Message: "containers.hello.files.0.mode: Invalid type",
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
			Message: "containers.hello.files.0.mode: Invalid type",
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
			Message: "containers.hello.files.0: content is required",
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
			Message: "containers.hello.files.0.content: Invalid type",
		},
		{
			Name: "containers.*.files.*.content is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var file = hello["files"].([]interface{})[0].(map[string]interface{})
				file["content"] = []string{}
				return src
			}(),
			Message: "containers.hello.files.0.content: Array must have at least 1 items",
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
			Message: "containers.hello.volumes: Invalid type",
		},
		{
			Name: "containers.*.volumes is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["volumes"] = []interface{}{}
				return src
			}(),
			Message: "containers.hello.volumes: Array must have at least 1 items",
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
			Message: "containers.hello.volumes.0: source is required",
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
			Message: "containers.hello.volumes.0.source: Invalid type",
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
			Message: "containers.hello.volumes.0.source: Invalid type",
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
			Message: "containers.hello.volumes.0.path: Invalid type",
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
			Message: "containers.hello.volumes.0.path: Invalid type",
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
			Message: "containers.hello.volumes.0: target is required",
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
			Message: "containers.hello.volumes.0.target: Invalid type",
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
			Message: "containers.hello.volumes.0.target: Invalid type",
		},
		{
			Name: "containers.*.volumes.*.read_only is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["read_only"] = nil
				return src
			}(),
			Message: "containers.hello.volumes.0.read_only: Invalid type",
		},
		{
			Name: "containers.*.volumes.*.read_only is not a boolean",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				var volumes = hello["volumes"].([]interface{})[0].(map[string]interface{})
				volumes["read_only"] = 12
				return src
			}(),
			Message: "containers.hello.volumes.0.read_only: Invalid type",
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
			Message: "containers.hello.resources: Invalid type",
		},
		{
			Name: "containers.*.resources is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.resources: Must have at least 1 properties",
		},
		{
			Name: "containers.*.resources.limits is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["limits"] = nil
				return src
			}(),
			Message: "containers.hello.resources.limits: Invalid type",
		},
		{
			Name: "containers.*.resources.limits is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["limits"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.resources.limits: Must have at least 1 properties",
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
			Message: "containers.hello.resources.limits.memory: Invalid type",
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
			Message: "containers.hello.resources.limits.memory: Invalid type",
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
			Message: "containers.hello.resources.limits.memory: Invalid type",
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
			Message: "containers.hello.resources.limits.cpu: Invalid type",
		},
		{
			Name: "containers.*.resources.requests is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["requests"] = nil
				return src
			}(),
			Message: "containers.hello.resources.requests: Invalid type",
		},
		{
			Name: "containers.*.resources.requests is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["resources"].(map[string]interface{})["requests"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.resources.requests: Must have at least 1 properties",
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
			Message: "containers.hello.resources.requests.memory: Invalid type",
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
			Message: "containers.hello.resources.requests.memory: Invalid type",
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
			Message: "containers.hello.resources.requests.memory: Invalid type",
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
			Message: "containers.hello.resources.requests.cpu: Invalid type",
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
			Message: "containers.hello.livenessProbe: Invalid type",
		},
		{
			Name: "containers.*.livenessProbe is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.livenessProbe: Must have at least 1 properties",
		},
		{
			Name: "containers.*.livenessProbe.httpGet is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["livenessProbe"].(map[string]interface{})["httpGet"] = nil
				return src
			}(),
			Message: "containers.hello.livenessProbe.httpGet: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet: path is required",
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
			Message: "containers.hello.livenessProbe.httpGet.path: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.path: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.port: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.port: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders: Array must have at least 1 item",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders.0.name: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders.0.name: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders.0.value: Invalid type",
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
			Message: "containers.hello.livenessProbe.httpGet.httpHeaders.0.value: Invalid type",
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
			Message: "containers.hello.readinessProbe: Invalid type",
		},
		{
			Name: "containers.*.readinessProbe is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["readinessProbe"] = map[string]interface{}{}
				return src
			}(),
			Message: "containers.hello.readinessProbe: Must have at least 1 properties",
		},
		{
			Name: "containers.*.readinessProbe.httpGet is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var hello = src["containers"].(map[string]interface{})["hello"].(map[string]interface{})
				hello["readinessProbe"].(map[string]interface{})["httpGet"] = nil
				return src
			}(),
			Message: "containers.hello.readinessProbe.httpGet: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet: path is required",
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
			Message: "containers.hello.readinessProbe.httpGet.path: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.path: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.port: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.port: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders: Array must have at least 1 item",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders.0.name: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders.0.name: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders.0.value: Invalid type",
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
			Message: "containers.hello.readinessProbe.httpGet.httpHeaders.0.value: Invalid type",
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
			Message: "resources: Invalid type",
		},
		{
			Name: "resources is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["resources"] = map[string]interface{}{}
				return src
			}(),
			Message: "resources: Must have at least 1 properties",
		},
		{
			Name: "resources.* is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				src["resources"].(map[string]interface{})["db"] = nil
				return src
			}(),
			Message: "resources.db: Invalid type",
		},

		// resources.*.image
		//
		{
			Name: "resources.*.type is missing",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				delete(db, "type")
				return src
			}(),
			Message: "resources.db: type is required",
		},
		{
			Name: "resources.*.type is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["type"] = nil
				return src
			}(),
			Message: "resources.db.type: Invalid type",
		},
		{
			Name: "resources.*.type is not a sring",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["type"] = 12
				return src
			}(),
			Message: "resources.db.type: Invalid type",
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
			Message: "resources.db.metadata: Invalid type",
		},
		{
			Name: "resources.*.metadata is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["metadata"] = map[string]interface{}{}
				return src
			}(),
			Message: "resources.db.metadata: Must have at least 1 properties",
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
			Message: "resources.db.metadata.annotations: Invalid type",
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
			Message: "resources.db.metadata.annotations: Must have at least 1 properties",
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
			Message: "resources.db.metadata.annotations.one: Invalid type.",
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
			Message: "resources.db.metadata.annotations.one: Invalid type.",
		},

		// resources.*.properties
		//
		{
			Name: "resources.*.properties is not set",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["properties"] = nil
				return src
			}(),
			Message: "",
		},
		{
			Name: "resources.*.properties is empty",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["properties"] = map[string]interface{}{}
				return src
			}(),
			Message: "",
		},
		{
			Name: "resources.*.properties is ignored",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["properties"] = map[string]interface{}{
					"key": "value",
				}
				return src
			}(),
			Message: "",
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
			Message: "resources.db.params: Invalid type",
		},
		{
			Name: "resources.*.params is not an object",
			Src: func() map[string]interface{} {
				src := newTestDocument()
				var db = src["resources"].(map[string]interface{})["db"].(map[string]interface{})
				db["params"] = 12
				return src
			}(),
			Message: "resources.db.params: Invalid type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var src = gojsonschema.NewGoLoader(tt.Src)
			res, err := Validate(src)
			assert.NoError(t, err)

			if tt.Message == "" {
				assert.True(t, res.Valid())
			} else {
				assert.False(t, res.Valid())

				var errors = res.Errors()
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].String(), tt.Message)
			}
		})
	}
}
