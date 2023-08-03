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

package aip0121

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceMustSupportList = &lint.ServiceRule{
	Name: lint.NewRuleName(121, "resource-must-support-list"),
	LintService: func(s *desc.ServiceDescriptor) []lint.Problem {
		problems := []lint.Problem{}
		// Add the empty string to avoid adding multiple nil checks.
		// Just say it's valid instead.
		resourcesWithList := stringset.New("")
		resourcesWithOtherMethods := map[string]*desc.MessageDescriptor{}
		// Iterate all RPCs and try to find resources. Mark the
		// resources which have a List method, and which ones do not.
		for _, m := range s.GetMethods() {
			if utils.IsListMethod(m) {
				if msg := utils.GetListResourceMessage(m); msg != nil {
					t := utils.GetResource(msg).GetType()
					resourcesWithList.Add(t)
				}
			} else if utils.IsCreateMethod(m) || utils.IsUpdateMethod(m) || utils.IsGetMethod(m) {
				if msg := utils.GetResponseType(m); msg != nil {
					// Skip tracking Singleton resources, they do not need List.
					if utils.IsSingletonResource(msg) {
						continue
					}
					t := utils.GetResource(msg).GetType()
					resourcesWithOtherMethods[t] = msg
				}
			}
		}

		for t := range resourcesWithOtherMethods {
			if !resourcesWithList.Contains(t) {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Missing Standard List method for resource %q", t,
					),
					Descriptor: s,
				})
			}
		}
		return problems
	},
}
