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
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceMustSupportGet = &lint.ServiceRule{
	Name: lint.NewRuleName(121, "resource-must-support-get"),
	LintService: func(s protoreflect.ServiceDescriptor) []lint.Problem {
		var problems []lint.Problem
		var resourcesWithGet stringset.Set
		var resourcesWithOtherMethods stringset.Set

		// Iterate all RPCs and try to find resources. Mark the
		// resources which have a Get method, and which ones do not.
		for i := 0; i < s.Methods().Len(); i++ {
			m := s.Methods().Get(i)
			// Streaming methods do not count as standard methods even if they
			// look like them.
			if utils.IsStreaming(m) {
				continue
			}

			if utils.IsGetMethod(m) && utils.IsResource(utils.GetResponseType(m)) {
				t := utils.GetResource(m.Output()).GetType()
				resourcesWithGet.Add(t)
			} else if utils.IsCreateMethod(m) || utils.IsUpdateMethod(m) {
				if msg := utils.GetResponseType(m); msg != nil && utils.IsResource(msg) {
					t := utils.GetResource(msg).GetType()
					resourcesWithOtherMethods.Add(t)
				}
			} else if utils.IsListMethod(m) {
				if msg := utils.GetListResourceMessage(m); msg != nil && utils.IsResource(msg) {
					t := utils.GetResource(msg).GetType()
					resourcesWithOtherMethods.Add(t)
				}
			}
		}

		for t := range resourcesWithOtherMethods {
			if !resourcesWithGet.Contains(t) {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Missing Standard Get method for resource %q", t,
					),
					Descriptor: s,
				})
			}
		}
		return problems
	},
}
