// Copyright 2024 Google LLC
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

var returnPartialSuccessType = &lint.FieldRule{
	Name: lint.NewRuleName(217, "return-partial-success-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "return_partial_success"
	},
	LintField: func(f *desc.FieldDescriptor) (problems []lint.Problem) {
		if utils.GetTypeName(f) != "bool" {
			problems = append(problems, lint.Problem{
				Message:    "`return_partial_success` field should be a `bool`.",
				Suggestion: "bool",
				Descriptor: f,
				Location:   locations.FieldType(f),
			})
		}
		return
	},
}
