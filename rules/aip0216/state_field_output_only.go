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

package aip0216

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"strings"
)

var stateFieldOutputOnly = &lint.FieldRule{
	Name: lint.NewRuleName(216, "state-field-output-only"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		// We care about the name of the State enum type.
		// AIP 0216 makes no mention of the state field name.
		et := f.GetEnumType()
		if et == nil {
			return false
		}

		if !strings.HasSuffix(et.GetName(), "State") {
			return false
		}

		return true
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		behaviors := utils.GetFieldBehavior(f)
		if !behaviors.Contains(annotations.FieldBehavior_OUTPUT_ONLY.String()) {
			return []lint.Problem{
				{
					Message:    "state fields must have field_behavior OUTPUT_ONLY",
					Descriptor: f,
				},
			}
		}
		return nil
	},
}
