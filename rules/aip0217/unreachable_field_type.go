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

package aip0217

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var unreachableFieldType = &lint.FieldRule{
	Name: lint.NewRuleName(217, "unreachable-field-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "unreachable"
	},
	LintField: func(f *desc.FieldDescriptor) (problems []lint.Problem) {
		if !f.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    "unreachable field should be repeated.",
				Descriptor: f,
			})
		}
		if utils.GetTypeName(f) != "string" {
			problems = append(problems, lint.Problem{
				Message:    "unreachable field should be a string.",
				Suggestion: "string",
				Descriptor: f,
				Location:   locations.FieldType(f),
			})
		}
		return
	},
}
