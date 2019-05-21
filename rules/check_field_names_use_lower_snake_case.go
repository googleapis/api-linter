// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"fmt"
	"strings"

	pref "google.golang.org/protobuf/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	"github.com/stoewer/go-strcase"
)

func init() {
	registerRules(checkFieldNamesUseLowerSnakeCase())
}

// checkFieldNamesUseLowerSnakeCase returns a lint.Rule
// that checks if a field name is using lower_snake_case.
func checkFieldNamesUseLowerSnakeCase() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "naming_formats", "field_names"),
			Description:  "check that field names use lower snake case",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FieldCallback: func(d pref.FieldDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				fieldName := string(d.Name())
				suggestion := toLowerSnakeCase(fieldName)
				if fieldName != suggestion {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("field named %q should use lower_snake_case", fieldName),
						Suggestion: suggestion,
						Descriptor: d,
					})
				}
				return
			},
		},
	}
}

// toLowerSnakeCase converts s to lower_snake_case.
func toLowerSnakeCase(s string) string {
	return strings.ToLower(strcase.SnakeCase(s))
}
