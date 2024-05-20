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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	score "github.com/score-spec/score-go/types"
)

func mustLoadWorkload(t *testing.T, spec string) *score.Workload {
	t.Helper()
	var raw score.Workload
	require.NoError(t, yaml.Unmarshal([]byte(spec), &raw))
	return &raw
}

func mustAddWorkload(t *testing.T, s *State[NoExtras, NoExtras, NoExtras], spec string) *State[NoExtras, NoExtras, NoExtras] {
	t.Helper()
	w := mustLoadWorkload(t, spec)
	n, err := s.WithWorkload(w, nil, NoExtras{})
	require.NoError(t, err)
	return n
}

func checkAndResetGuids[s any](t *testing.T, resources map[ResourceUid]ScoreResourceState[s]) {
	for uid, s := range resources {
		assert.NotEmpty(t, s.Guid)
		s.Guid = "00000000-0000-0000-0000-000000000000"
		resources[uid] = s
	}
}

func TestWithWorkload(t *testing.T) {
	start := new(State[NoExtras, NoExtras, NoExtras])

	t.Run("one", func(t *testing.T) {
		next, err := start.WithWorkload(mustLoadWorkload(t, `
metadata:
  name: example
  annotations:
    acme.org/x: y
containers:
  hello-world:
    image: hi
resources:
  foo:
    type: thing
`), nil, NoExtras{})
		require.NoError(t, err)
		assert.Len(t, start.Workloads, 0)
		assert.Len(t, next.Workloads, 1)
		assert.Nil(t, next.Workloads["example"].File, nil)
		assert.Equal(t, score.Workload{
			Metadata: map[string]interface{}{
				"name": "example",
				"annotations": map[string]interface{}{
					"acme.org/x": "y",
				},
			},
			Containers: map[string]score.Container{"hello-world": {Image: "hi"}},
			Resources:  map[string]score.Resource{"foo": {Type: "thing"}},
		}, next.Workloads["example"].Spec)
	})

	t.Run("two", func(t *testing.T) {
		next1, err := start.WithWorkload(mustLoadWorkload(t, `
metadata:
  name: example1
containers:
  hello-world:
    image: hi
resources:
  foo:
    type: thing
`), nil, NoExtras{})
		require.NoError(t, err)
		next2, err := next1.WithWorkload(mustLoadWorkload(t, `
metadata:
  name: example2
containers:
  hello-world:
    image: hi
`), nil, NoExtras{})
		require.NoError(t, err)

		assert.Len(t, start.Workloads, 0)
		assert.Len(t, next1.Workloads, 1)
		assert.Len(t, next2.Workloads, 2)
	})
}

func TestWithPrimedResources(t *testing.T) {
	start := new(State[NoExtras, NoExtras, NoExtras])

	t.Run("empty", func(t *testing.T) {
		next, err := start.WithPrimedResources()
		require.NoError(t, err)
		assert.Len(t, next.Resources, 0)
	})

	t.Run("one workload - nominal", func(t *testing.T) {
		next := mustAddWorkload(t, start, `
metadata: {"name": "example"}
resources:
  one:
    type: thing
  two:
    type: thing2
    class: banana
  three:
    type: thing3
    class: apple
    id: dog
    metadata:
      annotations:
        foo: bar
    params:
      color: green
  four:
    type: thing4
    id: elephant
  five:
    type: thing4
    id: elephant
    metadata:
      x: y
    params:
      color: blue
`)
		next, err := next.WithPrimedResources()
		require.NoError(t, err)
		assert.Len(t, start.Resources, 0)
		checkAndResetGuids(t, next.Resources)

		assert.Equal(t, map[ResourceUid]ScoreResourceState[NoExtras]{
			"thing.default#example.one": {
				Guid: "00000000-0000-0000-0000-000000000000",
				Type: "thing", Class: "default", Id: "example.one", State: map[string]interface{}{}, Outputs: map[string]interface{}{},
				SourceWorkload: "example",
			},
			"thing2.banana#example.two": {
				Guid: "00000000-0000-0000-0000-000000000000",
				Type: "thing2", Class: "banana", Id: "example.two", State: map[string]interface{}{}, Outputs: map[string]interface{}{},
				SourceWorkload: "example",
			},
			"thing3.apple#dog": {
				Guid: "00000000-0000-0000-0000-000000000000",
				Type: "thing3", Class: "apple", Id: "dog", State: map[string]interface{}{},
				Metadata:       map[string]interface{}{"annotations": map[string]interface{}{"foo": "bar"}},
				Params:         map[string]interface{}{"color": "green"},
				SourceWorkload: "example",
				Outputs:        map[string]interface{}{},
			},
			"thing4.default#elephant": {
				Guid: "00000000-0000-0000-0000-000000000000",
				Type: "thing4", Class: "default", Id: "elephant", State: map[string]interface{}{},
				Metadata:       map[string]interface{}{"x": "y"},
				Params:         map[string]interface{}{"color": "blue"},
				SourceWorkload: "example",
				Outputs:        map[string]interface{}{},
			},
		}, next.Resources)

		assert.NotNil(t, next.Resources["thing3.apple#dog"].Metadata["annotations"].(map[string]interface{}))
	})

	t.Run("one workload - same resource - same metadata", func(t *testing.T) {
		next := mustAddWorkload(t, start, `
metadata: {"name": "example"}
resources:
  one:
    type: thing
    id: elephant
    metadata:
      x: a
      annotations:
        acme.org/x: y
  two:
    type: thing
    id: elephant
    metadata:
      x: a
      annotations:
        acme.org/x: y
`)
		next, err := next.WithPrimedResources()
		assert.NoError(t, err)
		assert.Len(t, next.Resources, 1)
		assert.Equal(t, map[string]interface{}{
			"annotations": map[string]interface{}{
				"acme.org/x": "y",
			},
			"x": "a",
		}, next.Resources["thing.default#elephant"].Metadata)
	})

	t.Run("one workload - same resource - diff metadata", func(t *testing.T) {
		next := mustAddWorkload(t, start, `
metadata: {"name": "example"}
resources:
  one:
    type: thing
    id: elephant
    metadata:
      x: a
  two:
    type: thing
    id: elephant
    metadata:
      x: y
`)
		next, err := next.WithPrimedResources()
		require.EqualError(t, err, "resource 'thing.default#elephant': multiple definitions with different metadata")
		assert.Len(t, start.Resources, 0)
	})

	t.Run("one workload - diff params", func(t *testing.T) {
		next := mustAddWorkload(t, start, `
metadata: {"name": "example"}
resources:
  one:
    type: thing
    id: elephant
    params:
      x: a
  two:
    type: thing
    id: elephant
    params:
      x: y
`)
		next, err := next.WithPrimedResources()
		require.EqualError(t, err, "resource 'thing.default#elephant': multiple definitions with different params")
		assert.Len(t, start.Resources, 0)
	})

	t.Run("two workload - nominal", func(t *testing.T) {
		t.Run("one workload - nominal", func(t *testing.T) {
			next := mustAddWorkload(t, start, `
metadata: {"name": "example1"}
resources:
  one:
    type: thing
  two:
    type: thing2
    id: dog
`)
			next = mustAddWorkload(t, next, `
metadata: {"name": "example2"}
resources:
  one:
    type: thing
  two:
    type: thing2
    id: dog
`)
			next, err := next.WithPrimedResources()
			require.NoError(t, err)
			assert.Len(t, start.Resources, 0)
			assert.Len(t, next.Resources, 3)
			checkAndResetGuids(t, next.Resources)
			assert.Equal(t, map[ResourceUid]ScoreResourceState[NoExtras]{
				"thing.default#example1.one": {
					Guid: "00000000-0000-0000-0000-000000000000",
					Type: "thing", Class: "default", Id: "example1.one", State: map[string]interface{}{},
					SourceWorkload: "example1",
					Outputs:        map[string]interface{}{},
				},
				"thing.default#example2.one": {
					Guid: "00000000-0000-0000-0000-000000000000",
					Type: "thing", Class: "default", Id: "example2.one", State: map[string]interface{}{},
					SourceWorkload: "example2",
					Outputs:        map[string]interface{}{},
				},
				"thing2.default#dog": {
					Guid: "00000000-0000-0000-0000-000000000000",
					Type: "thing2", Class: "default", Id: "dog", State: map[string]interface{}{},
					SourceWorkload: "example1",
					Outputs:        map[string]interface{}{},
				},
			}, next.Resources)
		})
	})

}

func TestEncodeDecodeMetadata_yaml(t *testing.T) {
	state := new(State[NoExtras, NoExtras, NoExtras])
	state = mustAddWorkload(t, state, `
apiVersion: score.dev/v1b1
metadata:
  name: example
  annotations:
    acme.org/x: a
containers:
  main:
    image: nginx
resources:
  eg:
    type: example
    metadata:
      annotations:
        acme.org/x: b
`)
	raw, err := yaml.Marshal(state)
	assert.NoError(t, err)
	state = new(State[NoExtras, NoExtras, NoExtras])
	assert.NoError(t, yaml.Unmarshal(raw, &state))
	assert.Equal(t, score.WorkloadMetadata{"name": "example", "annotations": map[string]interface{}{"acme.org/x": "a"}}, state.Workloads["example"].Spec.Metadata)
	assert.Equal(t, score.ResourceMetadata{"annotations": map[string]interface{}{"acme.org/x": "b"}}, state.Workloads["example"].Spec.Resources["eg"].Metadata)

	state, err = state.WithPrimedResources()
	assert.NoError(t, err)
	assert.Equal(t, score.WorkloadMetadata{"name": "example", "annotations": map[string]interface{}{"acme.org/x": "a"}}, state.Workloads["example"].Spec.Metadata)
	assert.Equal(t, score.ResourceMetadata{"annotations": map[string]interface{}{"acme.org/x": "b"}}, state.Workloads["example"].Spec.Resources["eg"].Metadata)
	assert.Equal(t, map[string]interface{}{"annotations": map[string]interface{}{"acme.org/x": "b"}}, state.Resources["example.default#example.eg"].Metadata)
}

func TestEncodeDecodeMetadata_json(t *testing.T) {
	state := new(State[NoExtras, NoExtras, NoExtras])
	state = mustAddWorkload(t, state, `
apiVersion: score.dev/v1b1
metadata:
  name: example
  annotations:
    acme.org/x: a
containers:
  main:
    image: nginx
resources:
  eg:
    type: example
    metadata:
      annotations:
        acme.org/x: b
`)
	raw, err := json.Marshal(state)
	assert.NoError(t, err)
	state = new(State[NoExtras, NoExtras, NoExtras])
	assert.NoError(t, json.Unmarshal(raw, &state))
	assert.Equal(t, score.WorkloadMetadata{"name": "example", "annotations": map[string]interface{}{"acme.org/x": "a"}}, state.Workloads["example"].Spec.Metadata)
	assert.Equal(t, score.ResourceMetadata{"annotations": map[string]interface{}{"acme.org/x": "b"}}, state.Workloads["example"].Spec.Resources["eg"].Metadata)
}

func TestGetSortedResourceUids(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Empty(t, ru)
	})

	t.Run("one", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res": {Type: "thing", Params: map[string]interface{}{}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Equal(t, []ResourceUid{"thing.default#eg.res"}, ru)
	})

	t.Run("one cycle", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res.blah}"}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		_, err = s.GetSortedResourceUids()
		assert.EqualError(t, err, "a cycle exists involving resource param placeholders")
	})

	t.Run("two unrelated", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res1": {Type: "thing", Params: map[string]interface{}{}},
				"res2": {Type: "thing", Params: map[string]interface{}{}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Equal(t, []ResourceUid{"thing.default#eg.res1", "thing.default#eg.res2"}, ru)
	})

	t.Run("two linked", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res1": {Type: "thing", Params: map[string]interface{}{"x": "${resources.res2.blah}"}},
				"res2": {Type: "thing", Params: map[string]interface{}{}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Equal(t, []ResourceUid{"thing.default#eg.res2", "thing.default#eg.res1"}, ru)
	})

	t.Run("two cycle", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res1": {Type: "thing", Params: map[string]interface{}{"x": "${resources.res2.blah}"}},
				"res2": {Type: "thing", Params: map[string]interface{}{"y": "${resources.res1.blah}"}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		_, err = s.GetSortedResourceUids()
		assert.EqualError(t, err, "a cycle exists involving resource param placeholders")
	})

	t.Run("three linked", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res1": {Type: "thing", Params: map[string]interface{}{"x": "${resources.res2.blah}"}},
				"res2": {Type: "thing", Params: map[string]interface{}{}},
				"res3": {Type: "thing", Params: map[string]interface{}{"x": "${resources.res1.blah}"}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Equal(t, []ResourceUid{"thing.default#eg.res2", "thing.default#eg.res1", "thing.default#eg.res3"}, ru)
	})

	t.Run("complex", func(t *testing.T) {
		s, err := new(State[NoExtras, NoExtras, NoExtras]).WithWorkload(&score.Workload{
			Metadata: map[string]interface{}{"name": "eg"},
			Resources: map[string]score.Resource{
				"res1": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res2.blah} ${resources.res3.blah} ${resources.res4.blah} ${resources.res5.blah} ${resources.res6.blah}"}},
				"res2": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res3.blah} ${resources.res4.blah} ${resources.res5.blah} ${resources.res6.blah}"}},
				"res3": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res4.blah} ${resources.res5.blah} ${resources.res6.blah}"}},
				"res4": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res5.blah} ${resources.res6.blah}"}},
				"res5": {Type: "thing", Params: map[string]interface{}{"a": "${resources.res6.blah}"}},
				"res6": {Type: "thing", Params: map[string]interface{}{}},
			},
		}, nil, NoExtras{})
		assert.NoError(t, err)
		ru, err := s.GetSortedResourceUids()
		assert.NoError(t, err)
		assert.Equal(t, []ResourceUid{"thing.default#eg.res6", "thing.default#eg.res5", "thing.default#eg.res4", "thing.default#eg.res3", "thing.default#eg.res2", "thing.default#eg.res1"}, ru)
	})

}

type customStateExtras struct {
	Fruit string `yaml:"fruit"`
}

type customWorkloadExtras struct {
	Animal string `yaml:"animal"`
}

type customResourceExtras struct {
	Mineral string `yaml:"mineral"`
}

func TestCustomExtras(t *testing.T) {
	s := new(State[customStateExtras, customWorkloadExtras, customResourceExtras])
	s.Resources = map[ResourceUid]ScoreResourceState[customResourceExtras]{
		"thing.default#shared": {
			Guid: "00000000-0000-0000-0000-000000000000",
			Type: "thing", Class: "default", Id: "shared",
			Metadata: map[string]interface{}{},
			Params:   map[string]interface{}{},
			State:    map[string]interface{}{},
			Extras:   customResourceExtras{Mineral: "diamond"}},
	}
	s.SharedState = make(map[string]interface{})
	s.Extras.Fruit = "apple"
	s, _ = s.WithWorkload(&score.Workload{
		Metadata:   map[string]interface{}{"name": "eg"},
		Containers: map[string]score.Container{"example": {Image: "foo"}},
	}, nil, customWorkloadExtras{Animal: "bat"})

	raw, err := yaml.Marshal(s)
	assert.NoError(t, err)
	var rawOut map[string]interface{}
	assert.NoError(t, yaml.Unmarshal(raw, &rawOut))
	assert.Equal(t, map[string]interface{}{
		"workloads": map[string]interface{}{
			"eg": map[string]interface{}{
				"spec": map[string]interface{}{
					"apiVersion": "",
					"metadata":   map[string]interface{}{"name": "eg"},
					"containers": map[string]interface{}{
						"example": map[string]interface{}{
							"image": "foo",
						},
					},
				},
				"animal": "bat",
			},
		},
		"resources": map[string]interface{}{
			"thing.default#shared": map[string]interface{}{
				"guid":            "00000000-0000-0000-0000-000000000000",
				"type":            "thing",
				"class":           "default",
				"id":              "shared",
				"source_workload": "",
				"provisioner":     "",
				"mineral":         "diamond",
				"metadata":        map[string]interface{}{},
				"params":          map[string]interface{}{},
				"state":           map[string]interface{}{},
			},
		},
		"shared_state": map[string]interface{}{},
		"fruit":        "apple",
	}, rawOut)

	var s2 State[customStateExtras, customWorkloadExtras, customResourceExtras]
	assert.NoError(t, yaml.Unmarshal(raw, &s2))
	assert.Equal(t, "apple", s2.Extras.Fruit)
	assert.Equal(t, "bat", s2.Workloads["eg"].Extras.Animal)
	assert.Equal(t, "diamond", s2.Resources["thing.default#shared"].Extras.Mineral)
	assert.Equal(t, &s2, s)
}
