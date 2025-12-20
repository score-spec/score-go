// Copyright 2025 The Score Authors
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

package loader

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/score-spec/score-go/framework"
	"github.com/score-spec/score-go/types"
)

var (
	validplaceholderContent = regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
)

// ValidationError represets the set of non-schema validation issues with a
// workload.
type ValidationError struct {
	// Messages is the individual validation errors
	Messages []string `json:"messages"`
}

// Error returns a string representation of the error.
func (e *ValidationError) Error() string {
	return "validating workload:\n    " + strings.Join(e.Messages, "\n    ")
}

// allPlaceholdersInString returns all placeholders in the string.
// All plaecholders are returned, including duplicates.
func allPlaceholdersInString(s string) []string {
	placeholders := []string{}
	// SubstituteString only returns errors from the inner func or a missconfigured substitutor object
	framework.SubstituteString(s, func(placeholder string) (string, error) {
		placeholders = append(placeholders, placeholder)
		return "", nil
	})
	return placeholders
}

// allPlaceholdersIn returns all placeholders in the map or slice.
// All plaecholders are returned, including duplicates.
func allPlaceholdersIn(o any) []string {
	placeholders := []string{}
	// Substitute only returns errors from the inner func or a missconfigured substitutor object
	framework.Substitute(o, func(placeholder string) (string, error) {
		placeholders = append(placeholders, placeholder)
		return "", nil
	})
	return placeholders
}

// listAllPlaceholders returns all placeholders in the workload.
// The simplest parsing of placeholders is done:
// escape all $$ and then return the content of any ${...}
// The list is deduped, the order is undefined.
func listAllPlaceholders(workload *types.Workload) []string {
	placeholderSet := map[string]struct{}{}
	for _, container := range workload.Containers {
		for _, file := range container.Files {
			if (file.NoExpand == nil || !*file.NoExpand) && file.Content != nil {
				for _, placeholder := range allPlaceholdersInString(*file.Content) {
					placeholderSet[placeholder] = struct{}{}
				}
			}
		}
		for _, variable := range container.Variables {
			for _, placeholder := range allPlaceholdersInString(variable) {
				placeholderSet[placeholder] = struct{}{}
			}
		}
		for _, volume := range container.Volumes {
			for _, placeholder := range allPlaceholdersInString(volume.Source) {
				placeholderSet[placeholder] = struct{}{}
			}
		}
	}
	for _, resource := range workload.Resources {
		for _, placeholder := range allPlaceholdersIn(map[string]any(resource.Params)) {
			placeholderSet[placeholder] = struct{}{}
		}
	}
	placeholders := make([]string, len(placeholderSet))
	i := 0
	for k := range placeholderSet {
		placeholders[i] = k
		i++
	}
	return placeholders
}

// Validate checks for non-schame validation rules in the Score Spec.
//
// Validate returns multiple validation errors as a single
// ValidationError object. The individual messages can be extracted
// via the Messages property.
//
// The following validation rules are applied:
//
// - Placeholders must be well formed (contain at least two elements separated
// by ".", each element must be alphanumeric or contain "_" or "-")
//
// - The first element in a placeholder must be "resources" or "metadata"
//
// - All resource placeholders must resolve to a resource in the workload
func Validate(workload *types.Workload) error {
	errMsgs := []string{}
	placeholders := listAllPlaceholders(workload)
	for _, placeholder := range placeholders {
		if !validplaceholderContent.MatchString(placeholder) {
			errMsgs = append(errMsgs, fmt.Sprintf("placeholder ${%s} is malformed, must contain at least two elements separated by \".\", each element must be alphanumeric or contain \"_\" or \"-\"", placeholder))
			continue
		}
		// guaranteed to have at least 1 "." due to check above
		placeholderParts := strings.Split(placeholder, ".")
		switch placeholderParts[0] {
		case "resources":
			if workload.Resources != nil {
				if _, exists := workload.Resources[placeholderParts[1]]; !exists {
					errMsgs = append(errMsgs, fmt.Sprintf("placeholder ${%s} does not resolve to a resource, no resource with name \"%s\"", placeholder, placeholderParts[1]))
				}
			} else {
				errMsgs = append(errMsgs, fmt.Sprintf("placeholder ${%s} does not resolve to a resource, no resource with name \"%s\"", placeholder, placeholderParts[1]))
			}
		case "metadata":
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("placeholder ${%s} has unsupported first element of \"%s\"", placeholder, placeholderParts[0]))
		}
	}
	if len(errMsgs) > 0 {
		return &ValidationError{
			Messages: errMsgs,
		}
	}
	return nil
}
