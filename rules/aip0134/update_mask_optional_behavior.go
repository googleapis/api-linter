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

package aip0134

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var updateMaskOptionalBehavior = &lint.FieldRule{
	Name: lint.NewRuleName(134, "update-mask-optional-behavior"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "update_mask" && utils.IsUpdateRequestMessage(f.GetOwner())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		behaviors := utils.GetFieldBehavior(f)
		if !behaviors.Contains("OPTIONAL") {
			return []lint.Problem{
				{
					Message:    "Standard Update field `update_mask` must have `OPTIONAL` behavior",
					Descriptor: f,
					Location:   locations.FieldBehavior(f),
				},
			}
		}
		return nil
	},
}
