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

package aip0133

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Create methods should have a parent variable if the resource isn't top-level.
// This should be the only variable in the URI path.
var httpURIParent = &lint.MethodRule{
	Name: lint.NewRuleName(133, "http-uri-parent"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		// The response type of a Standard Create method must be the resource
		// itself, unless it is an LRO, in which case, the operation_info field
		// response_type must be the resource.
		res := m.Output()

		if info := utils.GetOperationInfo(m); info != nil {
			res = utils.FindMessage(m.ParentFile(), info.GetResponseType())
			// If we cannot resolve the response_type, then skip this check.
			// There is nothing we can do, and must let the check specific to
			// unresolvable operation_info types do its thing.
			if res == nil {
				return false
			}
		}

		return utils.IsCreateMethod(m) && utils.HasParent(utils.GetResource(res))
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		if problems := utils.LintHTTPURIHasParentVariable(m); problems != nil {
			return problems
		}
		if problems := utils.LintHTTPURIVariableCount(m, 1); problems != nil {
			return problems
		}
		return nil
	},
}
