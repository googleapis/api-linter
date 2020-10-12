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

package aip0154

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var fieldType = &lint.FieldRule{
	Name: lint.NewRuleName(154, "field-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "etag"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if t := utils.GetTypeName(f); t != "string" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The etag field should be a string, not %s.", t),
				Descriptor: f,
				Location:   locations.FieldType(f),
				Suggestion: "string",
			}}
		}
		return nil
	},
}
