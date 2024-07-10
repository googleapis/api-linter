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

// Package aip0123 contains rules defined in https://aip.dev/123.
package aip0123

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		123,
		duplicateResource,
		resourceAnnotation,
		resourceNameComponentsAlternate,
		resourceNameField,
		resourcePattern,
		resourcePlural,
		resourceReferenceType,
		resourceSingular,
		resourceTypeName,
		resourceVariables,
		resourceDefinitionVariables,
		resourceDefinitionPatterns,
		resourceDefinitionTypeName,
		resourcePatternSingular,
		nameNeverOptional,
	)
}

func isResourceMessage(m *desc.MessageDescriptor) bool {
	// If the parent of this message is a message, it is nested and shoudn't
	// be considered a resource, even if it has a name field.
	_, nested := m.GetParent().(*desc.MessageDescriptor)
	return m.FindFieldByName("name") != nil && !strings.HasSuffix(m.GetName(), "Request") &&
		!strings.HasSuffix(m.GetName(), "Response") && !nested
}

func hasResourceAnnotation(m *desc.MessageDescriptor) bool {
	return utils.GetResource(m) != nil
}

func hasResourceDefinitionAnnotation(f *desc.FileDescriptor) bool {
	return len(utils.GetResourceDefinitions(f)) > 0
}

// getVariables returns a slice of variables declared in the pattern.
//
// For example, a pattern of "publishers/{publisher}/books/{book}" would
// return []string{"publisher", "book"}.
func getVariables(pattern string) []string {
	answer := []string{}
	for _, match := range varRegexp.FindAllStringSubmatch(pattern, -1) {
		answer = append(answer, match[1])
	}
	return answer
}

// isRootLevelResourcePattern determines if the given pattern is that of a
// root-level resource by checking how many segments it has - root-level
// resource patterns have only two segments, thus one delimeter.
func isRootLevelResourcePattern(pattern string) bool {
	return strings.Count(pattern, "/") == 1
}

// getParentIDVariable is a helper that returns the parent resource ID segment
// for a given pattern. Returns empty string if the pattern has no variables or
// is a top-level resource.
func getParentIDVariable(pattern string) string {
	variables := getVariables(pattern)
	// no variables shouldn't happen but should be safe
	if len(variables) == 0 || isRootLevelResourcePattern(pattern) {
		return ""
	}

	// TODO: handle if singleton is a *parent*.
	if utils.IsSingletonResourcePattern(pattern) {
		// Last variable is the parent's for a singleton child.
		return variables[len(variables)-1]
	}

	return variables[len(variables)-2]
}

// nestedSingular returns the would be reduced singular form of a nested
// resource. Use isNestedName to check eligibility before using nestedSingular.
// This will return empty if the reosurce is not eligible for nested name
// reduction.
func nestedSingular(resource *apb.ResourceDescriptor) string {
	if !isNestedName(resource) {
		return ""
	}
	parentIDVar := getParentIDVariable(resource.GetPattern()[0])

	singular := utils.GetResourceSingular(resource)
	singularSnake := strcase.SnakeCase(singular)

	return strings.TrimPrefix(singularSnake, parentIDVar+"_")
}

// isNestedName determines if the resource naming could be reduced as a
// repetitive, nested collection as per AIP-122 Nested Collections. It does this
// by analyzing the first resource pattern defined, comparing the resource
// singular to the resource ID segment of the parent portion of the pattern. If
// that is a prefix of resource singular (in snake_case form as well), then the
// resource name could be reduced. For example, given `singular: "userEvent"`
// and `pattern: "users/{user}/userEvents/{user_event}"`, isNestedName would
// return `true`, because the `pattern` could be reduced to
// `"users/{user}/events/{event}"`.
func isNestedName(resource *apb.ResourceDescriptor) bool {
	if len(resource.GetPattern()) == 0 {
		return false
	}
	// only evaluate the first pattern b.c patterns cannot be reordered
	// and nested names must be used consistently, thus from the beginning
	pattern := resource.GetPattern()[0]

	// Can't be a nested collection if it is a top level resource.
	if isRootLevelResourcePattern(pattern) {
		return false
	}

	singular := utils.GetResourceSingular(resource)

	// If the resource type's singular is not camelCase then it is not a
	// multi-word type and we do not need to check for nestedness in naming.
	singularSnake := strcase.SnakeCase(singular)
	if strings.Count(singularSnake, "_") < 1 {
		return false
	}

	// Second to last resource ID variable will be the parent's to compare the
	// child resource's singular form against. This prevents us from needing to
	// deal with pluralizations due to the child resource typically being
	// pluralized on its own noun, not the parent noun.
	parentIDVar := getParentIDVariable(pattern)

	// For example:
	// publisher_credit starts with publisher --> nested
	// book_shelf does not start with library --> not nested re:naming
	return strings.HasPrefix(singularSnake, parentIDVar)
}

// getPlainPattern returns the pattern with all variables replaced with "*".
//
// For example, a pattern of "publishers/{publisher}/books/{book}" would
// return "publishers/*/books/*".
func getPlainPattern(pattern string) string {
	return varRegexp.ReplaceAllLiteralString(pattern, "*")
}

// getDesiredPattern returns the expected desired pattern, with errors we
// lint for corrected.
func getDesiredPattern(pattern string) string {
	want := []string{}
	for _, token := range strings.Split(pattern, "/") {
		if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
			varname := token[1 : len(token)-1]
			want = append(want, fmt.Sprintf("{%s}", strings.TrimSuffix(strcase.SnakeCase(varname), "_id")))
		} else {
			want = append(want, strcase.LowerCamelCase(token))
		}
	}
	return strings.Join(want, "/")
}

var varRegexp = regexp.MustCompile(`\{([^}=]+)}`)
