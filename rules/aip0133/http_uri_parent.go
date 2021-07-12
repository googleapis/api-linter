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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should have a parent variable if the resource isn't top-level.
// This should be the only variable in the URI path.
var httpURIParent = &lint.MethodRule{
	Name: lint.NewRuleName(133, "http-uri-parent"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return isCreateMethod(m) && !hasNoParent(m.GetOutputType())
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		if problems := utils.LintHTTPURIHasParentVariable(m); problems != nil {
			return problems
		}
		if problems := utils.LintHTTPURIVariableCount(m, 1); problems != nil {
			return problems
		}
		return nil
	},
}
