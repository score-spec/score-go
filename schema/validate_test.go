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

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestValidateYaml(t *testing.T) {
	var source = []byte(`
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
    class: external
  data:
    type: volume
    class: large
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

	err := ValidateYaml(bytes.NewReader(source))
	assert.NoError(t, err)
}

func TestValidateYaml_Error(t *testing.T) {
	var source = []byte(`
---
apiVersion: score.dev/v1b1
metadata:
  no-name: hello-world

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

	err := ValidateYaml(bytes.NewReader(source))
	assert.Error(t, err)
}

func TestValidateJson(t *testing.T) {
	var source = []byte(`
{
  "apiVersion": "score.dev/v1b1",
  "metadata": {
    "name": "hello-world"
  },
  "service": {
    "ports": {
      "www": {
        "port": 80,
        "targetPort": 8080
      }
    }
  },
  "containers": {
    "hello": {
      "image": "busybox",
      "command": [
        "/bin/echo"
      ],
      "args": [
        "Hello $(FRIEND)"
      ],
      "variables": {
        "FRIEND": "World!"
      },
      "files": [
        {
          "target": "/etc/hello-world/config.yaml",
          "mode": "666",
          "content": "${resources.env.APP_CONFIG}"
        }
      ],
      "volumes": [
        {
          "source": "${resources.data}",
          "path": "sub/path",
          "target": "/mnt/data",
          "readOnly": true
        }
      ],
      "resources": {
        "limits": {
          "memory": "128Mi",
          "cpu": "500m"
        },
        "requests": {
          "memory": "64Mi",
          "cpu": "250m"
        }
      },
      "livenessProbe": {
        "httpGet": {
          "path": "/alive",
          "port": 8080
        }
      },
      "readinessProbe": {
        "httpGet": {
          "path": "/ready",
          "port": 8080,
          "httpHeaders": [
            {
              "name": "Custom-Header",
              "value": "Awesome"
            }
          ]
        }
      }
    }
  },
  "resources": {
    "env": {
      "type": "environment"
    },
    "dns": {
      "type": "dns",
      "class": "external"
    },
    "data": {
      "type": "volume",
      "class": "large"
    },
    "db": {
      "type": "postgres",
      "metadata": {
        "annotations": {
          "my.org/version": "0.1"
        }
      },
      "params": {
        "extensions": {
          "uuid-ossp": {
            "schema": "uuid_schema",
            "version": "1.1"
          }
        }
      }
    }
  }
}
`)

	err := ValidateJson(bytes.NewReader(source))
	assert.NoError(t, err)
}

func TestValidateJson_Error(t *testing.T) {
	var source = []byte(`
{
  "apiVersion": "score.dev/v1b1",
  "metadata": {
    "no-name": "hello-world"
  },
  "service": {
    "ports": {
      "www": {
        "port": 80,
        "targetPort": 8080
      }
    }
  },
  "containers": {
    "hello": {
      "image": "busybox",
      "command": [
        "/bin/echo"
      ],
      "args": [
        "Hello $(FRIEND)"
      ],
      "variables": {
        "FRIEND": "World!"
      },
      "files": [
        {
          "target": "/etc/hello-world/config.yaml",
          "mode": "666",
          "content": "${resources.env.APP_CONFIG}"
        }
      ],
      "volumes": [
        {
          "source": "${resources.data}",
          "path": "sub/path",
          "target": "/mnt/data",
          "readOnly": true
        }
      ],
      "resources": {
        "limits": {
          "memory": "128Mi",
          "cpu": "500m"
        },
        "requests": {
          "memory": "64Mi",
          "cpu": "250m"
        }
      },
      "livenessProbe": {
        "httpGet": {
          "path": "/alive",
          "port": 8080
        }
      },
      "readinessProbe": {
        "httpGet": {
          "path": "/ready",
          "port": 8080,
          "httpHeaders": [
            {
              "name": "Custom-Header",
              "value": "Awesome"
            }
          ]
        }
      }
    }
  },
  "resources": {
    "env": {
      "type": "environment"
    },
    "dns": {
      "type": "dns"
    },
    "data": {
      "type": "volume"
    },
    "db": {
      "type": "postgres",
      "metadata": {
        "annotations": {
          "my.org/version": "0.1"
        }
      },
      "params": {
        "extensions": {
          "uuid-ossp": {
            "schema": "uuid_schema",
            "version": "1.1"
          }
        }
      }
    }
  }
}
`)

	err := ValidateJson(bytes.NewReader(source))
	assert.Error(t, err)
}

func TestValidateYaml_missing_workload_name(t *testing.T) {
	var source = []byte(`
---
apiVersion: score.dev/v1b1
metadata:
  no-name: hello-world
  something-else: value

containers:
  hello:
    image: busybox
`)

	err := ValidateYaml(bytes.NewReader(source))
	assert.EqualError(t, err, "jsonschema: '/metadata' does not validate with https://score.dev/schemas/score#/properties/metadata/required: missing properties: 'name'")
}

func TestApplyCommonUpgradeTransforms(t *testing.T) {
	var source = []byte(`
---
apiVersion: score.dev/v1b1
metadata:
  name: hello-world
containers:
  hello:
    image: busybox
    files:
    - target: /etc/hello-world/config.yaml
      mode: "666"
      content:
      - line1
      - line2
    volumes:
    - source: ${resources.data}
      target: /mnt/data
      read_only: true
`)

	var obj map[string]interface{}
	var dec = yaml.NewDecoder(bytes.NewReader(source))
	assert.NoError(t, dec.Decode(&obj))

	// first validation attempt should fail
	assert.Error(t, Validate(obj))

	// apply transforms
	changes, err := ApplyCommonUpgradeTransforms(obj)
	assert.NoError(t, err)
	assert.Len(t, changes, 2)

	// second validation attempt should succeed
	assert.NoError(t, Validate(obj))

	assert.Equal(t, "line1\nline2", obj["containers"].(map[string]interface{})["hello"].(map[string]interface{})["files"].([]interface{})[0].(map[string]interface{})["content"])
	assert.Equal(t, true, obj["containers"].(map[string]interface{})["hello"].(map[string]interface{})["volumes"].([]interface{})[0].(map[string]interface{})["readOnly"])
}
