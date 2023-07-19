// Copyright 2023 Google LLC
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

package aip0135

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Delete methods for resources that are parents should have a bool force field.
var forceField = &lint.MessageRule{
	Name: lint.NewRuleName(135, "force-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		name := m.FindFieldByName("name")
		ref := utils.GetResourceReference(name)
		validRef := ref != nil && ref.GetType() != "" && utils.FindResource(ref.GetType(), m.GetFile()) != nil

		return isDeleteRequestMessage(m) && validRef
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		force := m.FindFieldByName("force")
		name := m.FindFieldByName("name")
		ref := utils.GetResourceReference(name)
		res := utils.FindResource(ref.GetType(), m.GetFile())

		children := utils.FindResourceChildren(res, m.GetFile())
		if len(children) > 0 && force == nil {
			return []lint.Problem{
				{
					Message:    "Delete requests for resources with children should have a `bool force` field",
					Descriptor: m,
				},
			}
		}

		return nil
	},
}
