// Copyright 2024 Humanitec
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
	"crypto/rand"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"

	score "github.com/score-spec/score-go/types"
)

// State is the mega-structure that contains the state of our workload specifications and resources.
// Score specs are added to this structure and it stores the current resource set. Extra implementation specific fields
// are supported by the generic parameters.
type State[StateExtras any, WorkloadExtras any, ResourceExtras any] struct {
	Workloads   map[string]ScoreWorkloadState[WorkloadExtras]      `yaml:"workloads"`
	Resources   map[ResourceUid]ScoreResourceState[ResourceExtras] `yaml:"resources"`
	SharedState map[string]interface{}                             `yaml:"shared_state"`
	Extras      StateExtras                                        `yaml:",inline"`
}

// NoExtras can be used in place of the state or workload extras if no additional fields are needed.
type NoExtras struct {
}

// ScoreWorkloadState is the state stored per workload. We store the recorded workload spec, the file it came from if
// necessary to resolve relative references, and any extras for this implementation.
type ScoreWorkloadState[WorkloadExtras any] struct {
	// Spec is the final score spec after all overrides and images have been set. This is a validated score file.
	Spec score.Workload `yaml:"spec"`
	// File is the source score file if known.
	File *string `yaml:"file,omitempty"`
	// Extras stores any implementation specific extras needed for this workload.
	Extras WorkloadExtras `yaml:",inline"`
}

// ScoreResourceState is the state stored and tracked for each resource.
type ScoreResourceState[ResourceExtras any] struct {
	// Guid is a uuid assigned to this "instance" of the resource.
	Guid string `yaml:"guid"`
	// Type is the resource type.
	Type string `yaml:"type"`
	// Class is the resource class or 'default' if not provided.
	Class string `yaml:"class"`
	// Id is the generated id for the resource, either <workload>.<resName> or <id>. This is tracked so that
	// we can deduplicate and work out where a resource came from.
	Id string `yaml:"id"`

	Metadata map[string]interface{} `yaml:"metadata"`
	Params   map[string]interface{} `yaml:"params"`
	// SourceWorkload holds the workload name that had the best definition for this resource. "best" is either the
	// first one or the one with params defined.
	SourceWorkload string `yaml:"source_workload"`

	// ProvisionerUri is the resolved provisioner uri that should be found in the config. This is tracked so that
	// we identify which provisioner was used for a particular instance of the resource.
	ProvisionerUri string `yaml:"provisioner"`
	// State is the internal state local to this resource. It will be persisted to disk when possible.
	State map[string]interface{} `yaml:"state"`

	// Outputs is the current set of outputs for the resource. This is the output of calling the provider. It may contain
	// secrets so be careful when persisting this to disk.
	Outputs map[string]interface{} `yaml:"outputs,omitempty"`
	// OutputLookupFunc is function that allows certain in-process providers to defer any output generation. If this is
	// not provided, it will fall back to using what's in the outputs.
	OutputLookupFunc OutputLookupFunc `yaml:"-"`

	// Extras stores any implementation specific extras needed for this resource.
	Extras ResourceExtras `yaml:",inline"`
}

type OutputLookupFunc func(keys ...string) (interface{}, error)

// WithWorkload returns a new copy of State with the workload added, if the workload already exists with the same name
// then it will be replaced.
// This is not a deep copy, but any writes are executed in a copy-on-write manner to avoid modifying the source.
func (s *State[StateExtras, WorkloadExtras, ResourceExtras]) WithWorkload(spec *score.Workload, filePath *string, extras WorkloadExtras) (*State[StateExtras, WorkloadExtras, ResourceExtras], error) {
	out := *s
	if s.Workloads == nil {
		out.Workloads = make(map[string]ScoreWorkloadState[WorkloadExtras])
	} else {
		out.Workloads = maps.Clone(s.Workloads)
	}

	name, ok := spec.Metadata["name"].(string)
	if !ok {
		return nil, fmt.Errorf("metadata: name: is missing or is not a string")
	}
	out.Workloads[name] = ScoreWorkloadState[WorkloadExtras]{
		Spec:   *spec,
		File:   filePath,
		Extras: extras,
	}
	return &out, nil
}

// uuidV4 generates a uuid v4 string without dependencies
func uuidV4() string {
	// read 16 random bytes
	d := make([]byte, 16)
	_, _ = rand.Read(d)
	// set the version to version 4 (the top 4 bits of the 7th byte)
	d[6] = (d[6] & 0b_0000_1111) | 0b_0100_0000
	// format and print the output
	return fmt.Sprintf("%x-%x-%x-%x-%x", d[:4], d[4:6], d[6:8], d[8:10], d[10:])
}

func sortedStringMapKeys[v any](input map[string]v) []string {
	out := make([]string, 0, len(input))
	for s := range input {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// WithPrimedResources returns a new copy of State with all workload resources resolved to at least their initial type,
// class and id. New resources will have an empty provider set. Existing resources will not be touched.
// This is not a deep copy, but any writes are executed in a copy-on-write manner to avoid modifying the source.
func (s *State[StateExtras, WorkloadExtras, ResourceExtras]) WithPrimedResources() (*State[StateExtras, WorkloadExtras, ResourceExtras], error) {
	out := *s
	if s.Resources == nil {
		out.Resources = make(map[ResourceUid]ScoreResourceState[ResourceExtras])
	} else {
		out.Resources = maps.Clone(s.Resources)
	}

	primedResourceUids := make(map[ResourceUid]bool)
	for _, workloadName := range sortedStringMapKeys(s.Workloads) {
		workload := s.Workloads[workloadName]
		for _, resName := range sortedStringMapKeys(workload.Spec.Resources) {
			res := workload.Spec.Resources[resName]
			resUid := NewResourceUid(workloadName, resName, res.Type, res.Class, res.Id)
			if existing, ok := out.Resources[resUid]; !ok {
				out.Resources[resUid] = ScoreResourceState[ResourceExtras]{
					Guid:           uuidV4(),
					Type:           resUid.Type(),
					Class:          resUid.Class(),
					Id:             resUid.Id(),
					Metadata:       res.Metadata,
					Params:         res.Params,
					SourceWorkload: workloadName,
					State:          map[string]interface{}{},
					Outputs:        map[string]interface{}{},
				}
				primedResourceUids[resUid] = true
			} else if !primedResourceUids[resUid] {
				existing.Metadata = res.Metadata
				existing.Params = res.Params
				existing.SourceWorkload = workloadName
				out.Resources[resUid] = existing
				primedResourceUids[resUid] = true
			} else {
				// multiple definitions of the same shared resource, let's check for conflicting params and metadata
				if res.Params != nil {
					if existing.Params != nil && !reflect.DeepEqual(existing.Params, map[string]interface{}(res.Params)) {
						return nil, fmt.Errorf("resource '%s': multiple definitions with different params", resUid)
					}
					existing.Params = res.Params
					existing.SourceWorkload = workloadName
				}
				if res.Metadata != nil {
					if existing.Metadata != nil && !reflect.DeepEqual(existing.Metadata, map[string]interface{}(res.Metadata)) {
						return nil, fmt.Errorf("resource '%s': multiple definitions with different metadata", resUid)
					}
					existing.Metadata = res.Metadata
				}
				out.Resources[resUid] = existing
			}
		}
	}
	return &out, nil
}

func (s *State[StateExtras, WorkloadExtras, ResourceExtras]) getResourceDependencies(workloadName, resName string) (map[ResourceUid]bool, error) {
	outMap := make(map[ResourceUid]bool)
	res := s.Workloads[workloadName].Spec.Resources[resName]
	if res.Params == nil {
		return nil, nil
	}
	_, err := Substitute((map[string]interface{})(res.Params), func(ref string) (string, error) {
		parts := SplitRefParts(ref)
		if len(parts) > 1 && parts[0] == "resources" {
			rr, ok := s.Workloads[workloadName].Spec.Resources[parts[1]]
			if ok {
				outMap[NewResourceUid(workloadName, parts[1], rr.Type, rr.Class, rr.Id)] = true
			} else {
				return ref, fmt.Errorf("refers to unknown resource names '%s'", parts[1])
			}
		}
		return ref, nil
	})
	if err != nil {
		return nil, fmt.Errorf("workload '%s' resource '%s': %w", workloadName, resName, err)
	}
	return outMap, nil
}

// GetSortedResourceUids returns a topological sorting of the resource uids. The output order is deterministic and
// ensures that any resource output placeholder statements are strictly evaluated after their referenced resource.
// If cycles are detected an error will be thrown.
func (s *State[StateExtras, WorkloadExtras, ResourceExtras]) GetSortedResourceUids() ([]ResourceUid, error) {

	// We're implementing Kahn's algorithm (https://en.wikipedia.org/wiki/Topological_sorting#Kahn's_algorithm).
	nodesWithNoIncomingEdges := make(map[ResourceUid]bool)
	incomingEdges := make(map[ResourceUid]map[ResourceUid]bool, len(s.Resources))

	// We must first gather all the dependencies of each resource. Many resources won't have dependencies and will go
	// straight into the no-incoming-edges set
	for workloadName, workload := range s.Workloads {
		for resName, res := range workload.Spec.Resources {
			deps, err := s.getResourceDependencies(workloadName, resName)
			if err != nil {
				return nil, err
			}
			resUid := NewResourceUid(workloadName, resName, res.Type, res.Class, res.Id)
			if len(deps) == 0 {
				nodesWithNoIncomingEdges[resUid] = true
			} else {
				incomingEdges[resUid] = deps
			}
		}
	}

	// set up the output list
	output := make([]ResourceUid, 0, len(nodesWithNoIncomingEdges)+len(incomingEdges))

	// now iterate through the nodes with no incoming edges and subtract them from the
	for len(nodesWithNoIncomingEdges) > 0 {

		// to get a stable set, we grab whatever is on the set and convert it to a sorted list
		subset := make([]ResourceUid, 0, len(nodesWithNoIncomingEdges))
		for uid := range nodesWithNoIncomingEdges {
			subset = append(subset, uid)
		}
		clear(nodesWithNoIncomingEdges)
		slices.Sort(subset)

		// we can bulk append the subset to the output
		output = append(output, subset...)

		// remove a node from the no-incoming edges set
		for _, fromUid := range subset {
			// now find any nodes that had an edge going from this node to them
			for toUid, m := range incomingEdges {
				if m[fromUid] {
					// and remove the edge
					delete(m, fromUid)
					// if there are no incoming edges, then move it to the no-incoming-edges set
					if len(m) == 0 {
						delete(incomingEdges, toUid)
						nodesWithNoIncomingEdges[toUid] = true
					}
				}
			}
		}
	}
	// if we make no progress then there are cycles
	if len(incomingEdges) > 0 {
		return nil, fmt.Errorf("a cycle exists involving resource param placeholders")
	}
	return output, nil
}

// GetResourceOutputForWorkload returns an output function per resource name in the given workload. This is for
// passing into the compose translation context to resolve placeholder references.
// This does not modify the state.
func (s *State[StateExtras, WorkloadExtras, ResourceExtras]) GetResourceOutputForWorkload(workloadName string) (map[string]OutputLookupFunc, error) {
	workload, ok := s.Workloads[workloadName]
	if !ok {
		return nil, fmt.Errorf("workload '%s': does not exist", workloadName)
	}
	out := make(map[string]OutputLookupFunc)

	for resName, res := range workload.Spec.Resources {
		resUid := NewResourceUid(workloadName, resName, res.Type, res.Class, res.Id)
		state, ok := s.Resources[resUid]
		if !ok {
			return nil, fmt.Errorf("workload '%s': resource '%s' (%s) is not primed", workloadName, resName, resUid)
		}
		out[resName] = state.OutputLookup
	}
	return out, nil
}

// OutputLookup is a function which can traverse an outputs tree to find a resulting key, this defers to the embedded
// output function if it exists.
func (s *ScoreResourceState[ResourceExtras]) OutputLookup(keys ...string) (interface{}, error) {
	if s.OutputLookupFunc != nil {
		return s.OutputLookupFunc(keys...)
	} else if len(keys) == 0 {
		return nil, fmt.Errorf("at least one lookup key is required")
	}
	var resolvedValue interface{}
	resolvedValue = s.Outputs
	for _, k := range keys {
		ok := resolvedValue != nil
		if ok {
			var mapV map[string]interface{}
			mapV, ok = resolvedValue.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("cannot lookup key '%s', context is not a map", k)
			}
			resolvedValue, ok = mapV[k]
		}
		if !ok {
			return "", fmt.Errorf("key '%s' not found", k)
		}
	}
	return resolvedValue, nil
}
