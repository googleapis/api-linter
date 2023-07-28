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

package aip0134

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var allowMissing = &lint.MessageRule{
	Name: lint.NewRuleName(134, "request-allow-missing-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		if !utils.IsUpdateRequestMessage(m) {
			return false
		}
		r := utils.DeclarativeFriendlyResource(m)
		return r != nil && !utils.IsSingletonResource(r)
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		for _, field := range m.GetFields() {
			if field.GetName() == "allow_missing" && utils.GetTypeName(field) == "bool" && !field.IsRepeated() {
				return nil
			}
		}
		return []lint.Problem{{
			Descriptor: m,
			Message:    "Update requests on declarative-friendly resources should include a singular `bool allow_missing` field.",
		}}
	},
}
