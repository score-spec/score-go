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

package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Ref[k any](in k) *k {
	return &in
}

func DerefOr[k any](in *k, def k) k {
	if in == nil {
		return def
	}
	return *in
}

func parseAndFormatResourceLimits(rl ResourcesLimits) string {
	c, m, err := ParseResourceLimits(rl)
	return fmt.Sprintf("%d %d %v", DerefOr(c, -1), DerefOr(m, -1), err)
}

func TestParseResourceLimits(t *testing.T) {
	assert.Equal(t, "-1 -1 <nil>", parseAndFormatResourceLimits(ResourcesLimits{}))
	assert.Equal(t, "1000 1000000 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("1"), Memory: Ref("1M")}))
	assert.Equal(t, "-1 -1 failed to parse cpus 'banana' as a number", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("banana"), Memory: nil}))
	assert.Equal(t, "-1 -1 failed to parse memory 'banana' as a number", parseAndFormatResourceLimits(ResourcesLimits{Cpu: nil, Memory: Ref("banana")}))
	assert.Equal(t, "200 128974848 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("200m"), Memory: Ref("123Mi")}))
}
