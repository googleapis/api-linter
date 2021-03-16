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

package aip0235

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// The Batch Delete request message should have parent field.
var requestParentField = &lint.MessageRule{
	Name: lint.NewRuleName(235, "request-parent-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		// Sanity check: If the resource has a pattern, and that pattern
		// contains only one variable, then a parent field is not expected.
		//
		// In order to parse out the pattern, we get the resource message
		// from the request name, then get the resource annotation from that,
		// and then inspect the pattern there (oy!).
		plural := strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Request"), "BatchDelete")
		singular := utils.ToSingular(plural)
		if msg := utils.FindMessage(m.GetFile(), singular); msg != nil {
			if resource := utils.GetResource(msg); resource != nil {
				for _, pattern := range resource.GetPattern() {
					if strings.Count(pattern, "{") == 1 {
						return false
					}
				}
			}
		}

		return isBatchDeleteRequestMessage(m)
	},
	LintMessage: utils.LintFieldPresentAndSingularString("parent"),
}
