/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
    metadata:
      annotations:
        "my.org/version": "0.1"
    type: postgres
    properties:
      host:
        type: string
        default: localhost
        required: true
      port:
        default: 5432
      user.name:
    params: {
      extensions: {
        uuid-ossp: {
          schema: "uuid_schema",
          version: "1.1"
        }
      }
    }
`)

	var err = ValidateYaml(source)
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
    metadata:
      annotations:
        "my.org/version": "0.1"
    type: postgres
    properties:
      host:
        type: string
        default: localhost
        required: true
      port:
        default: 5432
      user.name:
    params: {
      extensions: {
        uuid-ossp: {
          schema: "uuid_schema",
          version: "1.1"
        }
      }
    }
`)

	var err = ValidateYaml(source)
	assert.Error(t, ErrInvalid, err)
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
			"content": [
				"---",
				"${resources.env.APP_CONFIG}"
			]
			}
		],
		"volumes": [
			{
			"source": "${resources.data}",
			"path": "sub/path",
			"target": "/mnt/data",
			"read_only": true
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
		"type": "environment",
		"properties": {
			"APP_CONFIG": null
		}
		},
		"dns": {
		"type": "dns"
		},
		"data": {
		"type": "volume"
		},
		"db": {
		"metadata": {
			"annotations": {
			"my.org/version": "0.1"
			}
		},
		"type": "postgres",
		"properties": {
			"host": {
			"type": "string",
			"default": "localhost",
			"required": true
			},
			"port": {
			"default": 5432
			},
			"user.name": null
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

	var err = ValidateJson(source)
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
			"content": [
				"---",
				"${resources.env.APP_CONFIG}"
			]
			}
		],
		"volumes": [
			{
			"source": "${resources.data}",
			"path": "sub/path",
			"target": "/mnt/data",
			"read_only": true
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
		"type": "environment",
		"properties": {
			"APP_CONFIG": null
		}
		},
		"dns": {
		"type": "dns"
		},
		"data": {
		"type": "volume"
		},
		"db": {
		"metadata": {
			"annotations": {
			"my.org/version": "0.1"
			}
		},
		"type": "postgres",
		"properties": {
			"host": {
			"type": "string",
			"default": "localhost",
			"required": true
			},
			"port": {
			"default": 5432
			},
			"user.name": null
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

	var err = ValidateJson(source)
	assert.Error(t, ErrInvalid, err)
}
