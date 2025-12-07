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

	"github.com/score-spec/score-go/types"
)

var (
	// Differnete from framework version as does not match escaped placeholders
	placeholderMatch        = regexp.MustCompile(`\$\{[^}]+\}`)
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
	// Escaping rule for $ is every pair of $$ should bcome a literal $
	// This means that $$${placeholder} should resolve to $<placeholder value>.
	// As we are only interested in the placeholder and not the rest of the
	// string, we can replace all $$ with another character and then look
	// for the placeholder.
	return placeholderMatch.FindAllString(strings.ReplaceAll(s, "$$", "_"), -1)
}

// allPlaceholdersIn returns all placeholders in the map or slice.
// All plaecholders are returned, including duplicates.
func allPlaceholdersIn(o any) []string {
	placeholders := []string{}
	if o == nil {
		return placeholders
	}
	switch v := o.(type) {
	case map[string]any:
		for _, val := range v {
			placeholders = append(placeholders, allPlaceholdersIn(val)...)
		}
	case []any:
		for _, val := range v {
			placeholders = append(placeholders, allPlaceholdersIn(val)...)
		}
	case string:
		placeholders = allPlaceholdersInString(v)
	}
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
					placeholderSet[placeholder[2:len(placeholder)-1]] = struct{}{}
				}
			}
		}
		for _, variable := range container.Variables {
			for _, placeholder := range allPlaceholdersInString(variable) {
				placeholderSet[placeholder[2:len(placeholder)-1]] = struct{}{}
			}
		}
		for _, volume := range container.Volumes {
			for _, placeholder := range allPlaceholdersInString(volume.Source) {
				placeholderSet[placeholder[2:len(placeholder)-1]] = struct{}{}
			}
		}
	}
	for _, resource := range workload.Resources {
		for _, placeholder := range allPlaceholdersIn(map[string]any(resource.Params)) {
			placeholderSet[placeholder[2:len(placeholder)-1]] = struct{}{}
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
			errMsgs = append(errMsgs, fmt.Sprintf("placeholder ${%s} has unknown first element of \"%s\"", placeholder, placeholderParts[0]))
		}
	}
	if len(errMsgs) > 0 {
		return &ValidationError{
			Messages: errMsgs,
		}
	}
	return nil
}
