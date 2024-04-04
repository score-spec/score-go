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
	"fmt"
	"strings"
)

// ResourceUid is a string for a unique resource identifier. This must be constructed through NewResourceUid
type ResourceUid string

// NewResourceUid constructs a new ResourceUid string.
func NewResourceUid(workloadName string, resName string, resType string, resClass *string, resId *string) ResourceUid {
	if resClass == nil {
		defaultClass := "default"
		resClass = &defaultClass
	}
	if resId != nil {
		return ResourceUid(fmt.Sprintf("%s.%s#%s", resType, *resClass, *resId))
	}
	return ResourceUid(fmt.Sprintf("%s.%s#%s.%s", resType, *resClass, workloadName, resName))
}

// Type returns the type of the resource
func (r ResourceUid) Type() string {
	return string(r)[0:strings.Index(string(r), ".")]
}

// Class returns the class of the resource, defaulted to "default"
func (r ResourceUid) Class() string {
	return string(r)[strings.Index(string(r), ".")+1 : strings.Index(string(r), "#")]
}

// Id returns the id of the resource, either <workload>.<name> or <id> if id is specified
func (r ResourceUid) Id() string {
	return string(r)[strings.Index(string(r), "#")+1:]
}
