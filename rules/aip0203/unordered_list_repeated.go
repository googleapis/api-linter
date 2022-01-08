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

package aip0203

import (
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var unorderedListRepeated = &lint.FieldRule{
	Name: lint.NewRuleName(203, "unordered-list-repeated"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.GetFieldBehavior(f).Contains("UNORDERED_LIST")
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.IsRepeated() {
			return nil
		}
		return []lint.Problem{{
			Message:    "The UNORDERED_LIST `google.api.field_behavior` annotation must not be applied to non-repeated fields.",
			Descriptor: f,
		}}
	},
}
