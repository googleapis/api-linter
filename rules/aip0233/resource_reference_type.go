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

package aip0233

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// BatchCreate methods should reference the target resource via `child_type` or
// the parent directly via `type`.
var resourceReferenceType = &lint.MethodRule{
	Name: lint.NewRuleName(233, "resource-reference-type"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		// Return type of the RPC.
		ot := utils.GetResponseType(m)
		if ot == nil {
			return false
		}

		// First repeated message field must be annotated with google.api.resource.
		repeated := utils.GetRepeatedMessageFields(ot)
		if len(repeated) == 0 {
			return false
		}
		resMsg := repeated[0].Message()
		resource := utils.GetResource(resMsg)

		parent := m.Input().Fields().ByName("parent")
		return isBatchCreateMethod(m) && parent != nil && utils.GetResourceReference(parent) != nil && resource != nil
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		// Return type of the RPC.
		ot := utils.GetResponseType(m)
		repeated := utils.GetRepeatedMessageFields(ot)
		resMsg := repeated[0].Message()

		resource := utils.GetResource(resMsg)
		parent := m.Input().Fields().ByName("parent")
		ref := utils.GetResourceReference(parent)

		if resource.GetType() == ref.GetType() {
			return []lint.Problem{{
				Message:    "BatchCreate should use a `child_type` reference to the created resource, not a `type` reference.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}
		if ref.GetChildType() != "" && resource.GetType() != ref.GetChildType() {
			return []lint.Problem{{
				Message:    "BatchCreate should use a `child_type` reference to the created resource.",
				Descriptor: parent,
				Location:   locations.FieldResourceReference(parent),
			}}
		}

		return nil
	},
}
