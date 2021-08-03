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

// LintFieldPresent returns a problem if the given message does not have the given field.
func LintFieldPresent(m *desc.MessageDescriptor, field string) (*desc.FieldDescriptor, []lint.Problem) {
	f := m.FindFieldByName(field)
	if f == nil {
		return nil, []lint.Problem{{
			Message:    fmt.Sprintf("Message `%s` has no `%s` field.", m.GetName(), field),
			Descriptor: m,
		}}
	}
	return f, nil
}

// LintSingularStringField returns a problem if the field is not a singular string.
func LintSingularStringField(f *desc.FieldDescriptor) []lint.Problem {
	return lintSingularField(f, builder.FieldTypeString(), "string")
}

func lintSingularField(f *desc.FieldDescriptor, t *builder.FieldType, want string) []lint.Problem {
	if f.GetType() != t.GetType() || f.IsRepeated() {
		return []lint.Problem{{
			Message:    fmt.Sprintf("The `%s` field must be a singular %s.", f.GetName(), want),
			Suggestion: want,
			Descriptor: f,
			Location:   locations.FieldType(f),
		}}
	}
	return nil
}

// LintSingularBoolField returns a problem if the field is not a singular bool.
func LintSingularBoolField(f *desc.FieldDescriptor) []lint.Problem {
	return lintSingularField(f, builder.FieldTypeBool(), "bool")
}

// LintFieldPresentAndSingularString returns a problem if a message does not have the given singular-string field.
func LintFieldPresentAndSingularString(field string) func(*desc.MessageDescriptor) []lint.Problem {
	return func(m *desc.MessageDescriptor) []lint.Problem {
		f, problems := LintFieldPresent(m, field)
		if f == nil {
			return problems
		}
		return LintSingularStringField(f)
	}
}

func lintFieldBehavior(f *desc.FieldDescriptor, want string) []lint.Problem {
	if !GetFieldBehavior(f).Contains(want) {
		return []lint.Problem{{
			Message:    fmt.Sprintf("The `%s` field should include `(google.api.field_behavior) = %s`.", f.GetName(), want),
			Descriptor: f,
		}}
	}
	return nil
}

// LintRequiredField returns a problem if the field's behavior is not REQUIRED.
func LintRequiredField(f *desc.FieldDescriptor) []lint.Problem {
	return lintFieldBehavior(f, "REQUIRED")
}

// LintOutputOnlyField returns a problem if the field's behavior is not OUTPUT_ONLY.
func LintOutputOnlyField(f *desc.FieldDescriptor) []lint.Problem {
	return lintFieldBehavior(f, "OUTPUT_ONLY")
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

func lintHTTPBody(m *desc.MethodDescriptor, want, msg string) []lint.Problem {
	for _, httpRule := range GetHTTPRules(m) {
		if httpRule.Body != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The `%s` method should %s HTTP body.", m.GetName(), msg),
				Descriptor: m,
				Location:   locations.MethodHTTPRule(m),
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
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}
		return nil
	}
}

// LintMethodHasMatchingRequestName returns a problem if the given method's request type does not
// have a name matching the method's, with a "Request" suffix.
func LintMethodHasMatchingRequestName(m *desc.MethodDescriptor) []lint.Problem {
	if got, want := m.GetInputType().GetName(), m.GetName()+"Request"; got != want {
		return []lint.Problem{{
			Message:    fmt.Sprintf("Request message should be named after the RPC, i.e. %q.", want),
			Suggestion: want,
			Descriptor: m,
			Location:   locations.MethodRequestType(m),
		}}
	}
	return nil
}

// LintMethodHasMatchingResponseName returns a problem if the given method's response type does not
// have a name matching the method's, with a "Response" suffix.
func LintMethodHasMatchingResponseName(m *desc.MethodDescriptor) []lint.Problem {
	if got, want := m.GetOutputType().GetName(), m.GetName()+"Response"; got != want {
		return []lint.Problem{{
			Message:    fmt.Sprintf("Response message should be named after the RPC, i.e. %q.", want),
			Suggestion: want,
			Descriptor: m,
			Location:   locations.MethodResponseType(m),
		}}
	}
	return nil
}

// LintHTTPURIHasParentVariable returns a problem if any of the given method's HTTP rules do not
// have a parent variable in the URI.
func LintHTTPURIHasParentVariable(m *desc.MethodDescriptor) []lint.Problem {
	return LintHTTPURIHasVariable(m, "parent")
}

// LintHTTPURIHasVariable returns a problem if any of the given method's HTTP rules do not
// have the given variable in the URI.
func LintHTTPURIHasVariable(m *desc.MethodDescriptor, v string) []lint.Problem {
	for _, httpRule := range GetHTTPRules(m) {
		if _, ok := httpRule.GetVariables()[v]; !ok {
			return []lint.Problem{{
				Message:    fmt.Sprintf("HTTP URI should include a `%s` variable.", v),
				Descriptor: m,
				Location:   locations.MethodHTTPRule(m),
			}}
		}
	}
	return nil
}

// LintHTTPURIVariableCount returns a problem if the given method's HTTP rules
// do not contain the given number of variables in the URI.
func LintHTTPURIVariableCount(m *desc.MethodDescriptor, n int) []lint.Problem {
	varsText := "variables"
	if n == 1 {
		varsText = "variable"
	}

	varsCount := 0
	for _, httpRule := range GetHTTPRules(m) {
		varsCount = max(varsCount, len(httpRule.GetVariables()))
	}
	if varsCount != n {
		return []lint.Problem{{
			Message:    fmt.Sprintf("HTTP URI should contain %d %s.", n, varsText),
			Descriptor: m,
			Location:   locations.MethodHTTPRule(m),
		}}
	}
	return nil
}

// LintHTTPURIHasNameVariable returns a problem if any of the given method's HTTP rules do not
// have a name variable in the URI.
func LintHTTPURIHasNameVariable(m *desc.MethodDescriptor) []lint.Problem {
	return LintHTTPURIHasVariable(m, "name")
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
