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

package aip0165

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// The Purge request message should have parent field.
var requestParentField = &lint.MessageRule{
	Name: lint.NewRuleName(165, "request-parent-field"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		// Sanity check: If the resource has a pattern, and that pattern
		// contains no variables, then a parent field is not expected.
		//
		// In order to parse out the pattern, we get the resource message
		// from the response, then get the resource annotation from that,
		// and then inspect the pattern there (oy!).
		plural := strings.TrimPrefix(strings.TrimSuffix(m.Name(), "Request"), "Purge")
		if resp := utils.FindMessage(m.ParentFile(), fmt.Sprintf("Purge%sResponse", plural)); resp != nil {
			if paged := resp.FindFieldByName(strcase.SnakeCase(plural)); paged != nil {
				if resource := utils.GetResource(paged.GetMessageType()); resource != nil {
					for _, pattern := range resource.GetPattern() {
						if strings.Count(pattern, "{") == 0 {
							return false
						}
					}
				}
			}
		}

		return isPurgeRequestMessage(m)
	},
	LintMessage: utils.LintFieldPresentAndSingularString("parent"),
}
