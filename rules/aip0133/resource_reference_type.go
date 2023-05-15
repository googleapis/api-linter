// Copyright 2022 Google LLC
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

package aip0133

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should reference the target resource via `child_type` or the
// parent directly via `type`.
var resourceReferenceType = &lint.MethodRule{
	Name: lint.NewRuleName(133, "resource-reference-type"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		// Return type of the RPC.
		ot := m.GetOutputType()
		if ot.GetName() == "Operation" {
			ot = utils.GetOperationResponseType(m)
		}

		// Unresolvable response_type for an Operation results in nil here.
		resource := utils.GetResource(ot)
		p := m.GetInputType().FindFieldByName("parent")
		return utils.IsCreateMethod(m) && p != nil && utils.GetResourceReference(p) != nil && resource != nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Return type of the RPC.
		ot := m.GetOutputType()
		if ot.GetName() == "Operation" {
			ot = utils.GetOperationResponseType(m)
		}
		resource := utils.GetResource(ot)
		parent := m.GetInputType().FindFieldByName("parent")
		ref := utils.GetResourceReference(parent)

		if resource.GetType() == ref.GetType() {
			return []lint.Problem{{
				Message:    "Create should use a `child_type` reference to the created resource, not a `type` reference.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}
		if ref.GetChildType() != "" && resource.GetType() != ref.GetChildType() {
			return []lint.Problem{{
				Message:    "Create should use a `child_type` reference to the created resource.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}

		return nil
	},
}
