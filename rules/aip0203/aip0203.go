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
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AIP-203 rules to the provided registry.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		inputOnly,
	)
}

// check leading comment of a field and produce a problem
// if the comment matches the give pattern.
func checkAnnotation(f *desc.FieldDescriptor, pattern *regexp.Regexp) []lint.Problem {
	leadingComments := f.GetSourceInfo().GetLeadingComments()
	if pattern.MatchString(leadingComments) {
		return []lint.Problem{lint.Problem{
			Message:    fmt.Sprintf("Use the `google.api.field_behavior` annotation instead of %q.", pattern.FindString(leadingComments)),
			Descriptor: f,
		}}
	}
	return nil
}
