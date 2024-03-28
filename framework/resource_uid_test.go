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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceUid_basic(t *testing.T) {
	r := NewResourceUid("work", "my-res", "res-type", nil, nil)
	assert.Equal(t, "res-type.default#work.my-res", string(r))
	assert.Equal(t, "work.my-res", r.Id())
	assert.Equal(t, "res-type", r.Type())
	assert.Equal(t, "default", r.Class())
}

func TestResourceUid_with_class(t *testing.T) {
	someClass := "something"
	r := NewResourceUid("work", "my-res", "res-type", &someClass, nil)
	assert.Equal(t, "res-type.something#work.my-res", string(r))
	assert.Equal(t, "work.my-res", r.Id())
	assert.Equal(t, "res-type", r.Type())
	assert.Equal(t, "something", r.Class())
}

func TestResourceUid_with_id(t *testing.T) {
	someId := "something"
	r := NewResourceUid("work", "my-res", "res-type", nil, &someId)
	assert.Equal(t, "res-type.default#something", string(r))
	assert.Equal(t, "something", r.Id())
	assert.Equal(t, "res-type", r.Type())
	assert.Equal(t, "default", r.Class())
}
