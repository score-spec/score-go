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

package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubstitutionsResolve(t *testing.T) {
	substitutions := Substitutions{
		"resources.db.name":     {Value: "db_name"},
		"resources.db.port":     {Value: 5432},
		"resources.db.host":     {Value: "db.example.com"},
		"resources.db.username": {Secret: &SecretRef{Type: SecretTypeDirect, Value: "test-user"}},
		"resources.db.password": {Secret: &SecretRef{
			Type:  SecretTypeAWS,
			Store: "secrets",
			Ref:   "arn:aws:secretsmanager:us-east-1:0123456789:secret:creds/dbs/password",
		}},
		"resources.bucket": {Value: "my-bucket"},
	}
	testCases := []struct {
		name        string
		placeholder string
		expected    *SubValue
		shouldFail  bool
	}{
		{
			name:        "string value",
			placeholder: "resources.db.name",
			expected:    &SubValue{Value: "db_name"},
		},
		{
			name:        "non-string value",
			placeholder: "resources.db.port",
			expected:    &SubValue{Value: 5432},
		},
		{
			name:        "direct secret",
			placeholder: "resources.db.username",
			expected:    &SubValue{Secret: &SecretRef{Type: SecretTypeDirect, Value: "test-user"}},
		},
		{
			name:        "stored secret",
			placeholder: "resources.db.password",
			expected: &SubValue{Secret: &SecretRef{
				Type:  SecretTypeAWS,
				Store: "secrets",
				Ref:   "arn:aws:secretsmanager:us-east-1:0123456789:secret:creds/dbs/password",
			}},
		},
		{
			name:        "two parts",
			placeholder: "resources.bucket",
			expected:    &SubValue{Value: "my-bucket"},
		},
		{
			name:        "single part",
			placeholder: "resources",
			shouldFail:  true,
		},
		{
			name:        "wrapped in ${...}",
			placeholder: "${resources.db.name}",
			shouldFail:  true,
		},
		{
			name:        "name does not exist",
			placeholder: "resources.no-exist.name",
			shouldFail:  true,
		},
		// metadata and container placeholders are replaced when the spec is normalized.
		{
			name:        "metadata placeholder",
			placeholder: "metadata.name",
			shouldFail:  true,
		},
		{
			name:        "container placeholder",
			placeholder: "container.image",
			shouldFail:  true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := substitutions.Resolve(testCase.placeholder)
			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, actual)
			}
		})
	}
}

func TestSubstitutionsResolve_NoSubstitutions(t *testing.T) {
	var substitutions Substitutions

	t.Run("no metadata", func(t *testing.T) {
		_, err := substitutions.Resolve("metadata.name")
		assert.Error(t, err)
	})
	t.Run("no resources", func(t *testing.T) {
		_, err := substitutions.Resolve("resources.does-not-exist.name")
		assert.Error(t, err)
	})
	t.Run("no container", func(t *testing.T) {
		_, err := substitutions.Resolve("container.image")
		assert.Error(t, err)
	})
}
