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
	"fmt"
	"strings"
)

// SecretType identifies the secret store flavour and the format of a
// SecretRef's Ref. The list of predefined types follows the providers of the
// External Secrets Operator and is extensible: implementations are free to
// support other types.
type SecretType string

const (
	// SecretTypeDirect is used when the secret value is provided directly in
	// SecretRef.Value rather than stored in a secret store.
	SecretTypeDirect SecretType = "direct"
	// SecretTypeAWS is AWS Secrets Manager; Ref is the ARN of the secret.
	SecretTypeAWS SecretType = "aws"
	// SecretTypeAzureKV is Azure Key Vault; Ref is the URL of the key.
	SecretTypeAzureKV SecretType = "azurekv"
	// SecretTypeGCPSM is Google Secret Manager; Ref is the fully qualified
	// path starting with projects/{PROJECT_ID}/.
	SecretTypeGCPSM SecretType = "gcpsm"
	// SecretTypeKubernetes is a Kubernetes Secret; Ref is the K8s API URI of
	// the Secret.
	SecretTypeKubernetes SecretType = "kubernetes"
	// SecretTypeVault is HashiCorp Vault; Ref is the URL of the key.
	SecretTypeVault SecretType = "vault"
)

// SubValue represents a secret or non-secret value in the Substitutions map.
//
// Only one of Value or Secret can be specified.
// If both are nil, the value is assumed to be nil and not secret.
// (A secret nil value requires Secret to be non-nil.)
type SubValue struct {
	Value  any        `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
	Secret *SecretRef `json:"secret,omitempty" yaml:"secret,omitempty" mapstructure:"secret,omitempty"`
}

// SecretRef represents a secret value to be used in substitutions.
//
// It is a high-level abstraction of the approach used by the Kubernetes
// External Secrets Operator, hence Store and Ref. In addition, if secrets are
// dynamically managed and not persisted, they can be provided directly.
type SecretRef struct {
	// Type of secret reference being used.
	//
	// The type defines the secret store flavour and the format of the
	// reference. Required.
	Type SecretType `json:"type" yaml:"type" mapstructure:"type"`

	// Name of the secret store the secret is associated with.
	//
	// This can be used to associate credentials or access rights with the
	// secret reference. Not used when raw secret values are provided
	// (Type is SecretTypeDirect).
	// Optional.
	Store string `json:"store,omitempty" yaml:"store,omitempty" mapstructure:"store,omitempty"`

	// Identification of the secret in the secret store.
	//
	// It is opaque and interpreted according to Type. Where the store supports versions,
	// the version is encoded in Ref. If empty, Value must be supplied.
	Ref string `json:"ref,omitempty" yaml:"ref,omitempty" mapstructure:"ref,omitempty"`

	// If the secret does not reside in a store, the "plaintext" of the
	// secret can be provided. If Value is provided, Type should be SecretTypeDirect.
	// Type==SecretTypeDirect, Ref=="" and Value==nil implies a secret value of nil.
	Value any `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

// Resolve resolves a placeholder expression to a SubValue from the substitutions map.
//
// The placeholder is the expression inside ${...}, without the surrounding ${ and },
// e.g. "resources.db.host". It must contain two or more parts, separated with a dot.
// Only placeholders starting with "resources" are supported. Other types of placeholders
// should be replaced when the spec is normalized.
//
// The lookup is by the exact placeholder key: nested outputs and escaped dots are not
// interpreted, so the substitutions map must contain the full key.
func (s Substitutions) Resolve(placeholder string) (*SubValue, error) {
	parts := strings.Split(placeholder, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid placeholder: must have at least 2 parts, got %q", placeholder)
	}
	switch parts[0] {
	case "resources":
		if sub, exist := s[placeholder]; exist {
			return &sub, nil
		}
		return nil, fmt.Errorf("resolving placeholder %q: no substitution found", placeholder)
	default:
		return nil, fmt.Errorf("invalid placeholder: must start with \"resources\", got %q", placeholder)
	}
}
