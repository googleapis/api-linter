// Copyright 2020 Google LLC
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

package utils

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// LintStringField returns a problem if the field is not a string.
func LintStringField(f *desc.FieldDescriptor) []lint.Problem {
	if f.GetType() != builder.FieldTypeString().GetType() {
		return []lint.Problem{{
			Message:    fmt.Sprintf("The `%s` field must be a string.", f.GetType()),
			Suggestion: "string",
			Descriptor: f,
			Location:   locations.FieldType(f),
		}}
	}
	return nil
}

// LintRequiredField returns a problem if the field's behavior is not REQUIRED.
func LintRequiredField(f *desc.FieldDescriptor) []lint.Problem {
	if !GetFieldBehavior(f).Contains("REQUIRED") {
		return []lint.Problem{{
			Message:    fmt.Sprintf("The `%s` field should include `(google.api.field_behavior) = REQUIRED`.", f.GetName()),
			Descriptor: f,
		}}
	}
	return nil
}

// LintFieldResourceReference returns a problem if the field does not have a resource reference annotation.
func LintFieldResourceReference(f *desc.FieldDescriptor) []lint.Problem {
	if ref := GetResourceReference(f); ref == nil {
		return []lint.Problem{{
			Message:    fmt.Sprintf("The `%s` field should include a `google.api.resource_reference` annotation.", f.GetName()),
			Descriptor: f,
		}}
	}
	return nil
}

// LintParentFieldPresenceAndType returns a problem if the message does not have a `parent` field
// or if its type is not `string`.
func LintParentFieldPresenceAndType(m *desc.MessageDescriptor) []lint.Problem {
	parentField := m.FindFieldByName("parent")
	if parentField == nil {
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
			Descriptor: m,
		}}
	}
	return LintStringField(parentField)
}

func lintHTTPBody(m *desc.MethodDescriptor, want, msg string) []lint.Problem {
	for _, httpRule := range GetHTTPRules(m) {
		if httpRule.Body != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The `%s` method should %s HTTP body.", m.GetName(), msg),
				Descriptor: m,
			}}
		}
	}
	return nil
}

// LintNoHTTPBody returns a problem for each HTTP rule whose body is not "".
func LintNoHTTPBody(m *desc.MethodDescriptor) []lint.Problem {
	return lintHTTPBody(m, "", "not have an")
}

// LintWildcardHTTPBody returns a problem for each HTTP rule whose body is not "*".
func LintWildcardHTTPBody(m *desc.MethodDescriptor) []lint.Problem {
	return lintHTTPBody(m, "*", `use "*" as the`)
}

// LintHTTPMethod returns a problem for each HTTP rule whose HTTP method is not the given one.
func LintHTTPMethod(verb string) func(*desc.MethodDescriptor) []lint.Problem {
	return func(m *desc.MethodDescriptor) []lint.Problem {
		for _, httpRule := range GetHTTPRules(m) {
			if httpRule.Method != verb {
				return []lint.Problem{{
					Message:    fmt.Sprintf("The `%s` method should use the HTTP %s verb.", m.GetName(), verb),
					Descriptor: m,
				}}
			}
		}
		return nil
	}
}
