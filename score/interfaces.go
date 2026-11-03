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

// Package score defines the common interfaces and types shared by Score
// implementations, so that normalization, provisioning and generation logic
// can be reused, mixed and matched, and imported as a library.
package score

import (
	"context"
)

// Normalizer converts a Score file parsed from YAML into a normalized v1 spec.
//
// A normalized spec has all overrides applied, all images resolved, short-hand
// fully expanded, non-resource placeholders (e.g. workload, metadata,
// container) fully expanded and default values made explicit.
type Normalizer interface {
	// NormalizeToV1 converts rawSpec into a normalized v1 Workload, applying
	// overrides and resolving images.
	//
	// Handles all previous Score versions, upgrading as required.
	NormalizeToV1(rawSpec map[string]any, overrides map[string]any, images map[string]string) (*Workload, error)
}

// Provisioner provisions the resources required by a Score spec.
type Provisioner interface {
	// Provision provisions resources as required by the provided Score spec
	// and returns a substitution map for its resource outputs.
	Provision(ctx context.Context, spec *Workload) (Substitutions, error)
}

// Generator converts a Score spec into deployable artefacts for a target platform.
type Generator interface {
	// Generate produces deployable artefacts for the target platform, as a map
	// of names to data blobs.
	Generate(ctx context.Context, spec *Workload, subs Substitutions) (map[string][]byte, error)
}

// Substitutions maps placeholder keys (e.g. "resources.db.host") to their values.
type Substitutions map[string]SubValue
