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

package aip0132

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// List methods should reference the target resource via `child_type` or the
// parent directly via `type`.
var resourceReferenceType = &lint.MethodRule{
	Name: lint.NewRuleName(132, "resource-reference-type"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		p := m.GetInputType().FindFieldByName("parent")

		var resource *annotations.ResourceDescriptor
		resourceField := utils.GetListResourceMessage(m)
		if resourceField != nil {
			resource = utils.GetResource(resourceField)
		}
		return utils.IsListMethod(m) && p != nil && utils.GetResourceReference(p) != nil && resource != nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// The first repeated message field must be the paginated resource.
		repeated := utils.GetRepeatedMessageFields(m.GetOutputType())
		resource := utils.GetResource(repeated[0].GetMessageType())

		parent := m.GetInputType().FindFieldByName("parent")
		ref := utils.GetResourceReference(parent)

		if resource.GetType() == ref.GetType() {
			return []lint.Problem{{
				Message:    "List should use a `child_type` reference to the paginated resource, not a `type` reference.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}
		if ref.GetChildType() != "" && resource.GetType() != ref.GetChildType() {
			return []lint.Problem{{
				Message:    "List should use a `child_type` reference to the paginated resource.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}

		return nil
	},
}
