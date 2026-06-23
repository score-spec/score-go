// Copyright 2026 The Score Authors
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

package framework

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	score "github.com/score-spec/score-go/types"
)

func TestNetworkService(t *testing.T) {
	t.Run("round trip TCP", func(t *testing.T) {
		original := NetworkService{
			ServiceName: "hello",
			Ports: map[string]ServicePort{
				"http": {
					Name:       "http",
					Port:       8080,
					TargetPort: 80,
					Protocol:   score.ServicePortProtocolTCP,
				},
				"8080": {
					Name:       "http",
					Port:       8080,
					TargetPort: 80,
					Protocol:   score.ServicePortProtocolTCP,
				},
			},
		}

		raw, err := yaml.Marshal(original)
		require.NoError(t, err)

		var decoded NetworkService
		require.NoError(t, yaml.Unmarshal(raw, &decoded))
		assert.Equal(t, original, decoded)
	})

	t.Run("round trip UDP", func(t *testing.T) {
		original := NetworkService{
			ServiceName: "hello",
			Ports: map[string]ServicePort{
				"dns": {
					Name:       "dns",
					Port:       53,
					TargetPort: 5353,
					Protocol:   score.ServicePortProtocolUDP,
				},
			},
		}

		raw, err := yaml.Marshal(original)
		require.NoError(t, err)

		var decoded NetworkService
		require.NoError(t, yaml.Unmarshal(raw, &decoded))
		assert.Equal(t, original, decoded)
	})

	t.Run("yaml field names", func(t *testing.T) {
		input := `
service_name: my-service
ports:
  http:
    name: http
    port: 80
    target_port: 8080
    protocol: TCP
`

		var ns NetworkService
		require.NoError(t, yaml.Unmarshal([]byte(input), &ns))
		assert.Equal(t, "my-service", ns.ServiceName)
		assert.Equal(t, "http", ns.Ports["http"].Name)
		assert.Equal(t, 80, ns.Ports["http"].Port)
		assert.Equal(t, 8080, ns.Ports["http"].TargetPort)
		assert.Equal(t, score.ServicePortProtocol("TCP"), ns.Ports["http"].Protocol)
	})
}
