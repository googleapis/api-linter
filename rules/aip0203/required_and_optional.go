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

package aip0203

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var requiredAndOptional = &lint.FieldRule{
	Name: lint.NewRuleName(203, "required-and-optional"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.GetFieldBehavior(f).Contains("REQUIRED")
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.IsProto3Optional() {
			return []lint.Problem{{
				Message:    "Required fields must not use the optional keyword.",
				Descriptor: f,
				Location:   locations.FieldLabel(f),
				Suggestion: "",
			}}
		}
		return nil
	},
}
