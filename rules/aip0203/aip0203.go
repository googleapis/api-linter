// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package aip0203 contains rules defined in https://aip.dev/203.
package aip0203

import (
	"fmt"
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AIP-203 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		203,
		inputOnly,
		immutable,
		optional,
		optionalBehaviorConflict,
		optionalBehaviorConsistency,
		outputOnly,
		required,
		requiredAndOptional,
		unorderedListRepeated,
	)
}

// check leading comments of a field and produce a problem
// if the comments match the give pattern.
func checkLeadingComments(pattern *regexp.Regexp, annotation string, unless ...*regexp.Regexp) func(*desc.FieldDescriptor) []lint.Problem {
	return func(f *desc.FieldDescriptor) []lint.Problem {
		leadingComments := f.GetSourceInfo().GetLeadingComments()
		for _, ul := range unless {
			if ul.MatchString(leadingComments) {
				return nil
			}
		}
		if pattern.MatchString(leadingComments) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Use the `google.api.field_behavior` annotation instead of %q in the leading comments. For example, `string name = 1 [(google.api.field_behavior) = %s];`.", pattern.FindString(leadingComments), annotation),
				Descriptor: f,
			}}
		}
		return nil
	}
}

func withoutFieldBehavior(f *desc.FieldDescriptor) bool {
	return utils.GetFieldBehavior(f).Len() == 0
}
