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
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// TODO: test against LROs
var resourceMustSupportGet = &lint.ServiceRule{
	Name: lint.NewRuleName(121, "resource-must-support-get"),
	LintService: func(s *desc.ServiceDescriptor) []lint.Problem {
		problems := []lint.Problem{}
		// Add the empty string to avoid adding multiple nil checks.
		// Just say it's valid instead.
		resourcesWithGet := stringset.New("")
		resourcesWithOtherMethods := map[string]*desc.MessageDescriptor{}
		// Iterate all RPCs and try to find resources. Mark the
		// resources which have a Get method, and which ones do not.
		for _, m := range s.GetMethods() {
			if utils.IsGetMethod(m) {
				t := utils.GetResource(m.GetOutputType()).GetType()
				resourcesWithGet.Add(t)
			} else if utils.IsCreateMethod(m) || utils.IsUpdateMethod(m) {
				message := utils.GetResponseType(m)
				if message != nil {
					t := utils.GetResource(message).GetType()
					resourcesWithOtherMethods[t] = message
				}
			} else if utils.IsListMethod(m) {
				message := utils.GetListResourceMessage(m)
				if message != nil {
					t := utils.GetResource(message).GetType()
					resourcesWithOtherMethods[t] = message
				}
			}
		}

		for t, m := range resourcesWithOtherMethods {
			if !resourcesWithGet.Contains(t) {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Missing Standard Get method for resource %q", t,
					),
					Descriptor: m,
					Location:   locations.MessageResource(m),
				})
			}
		}
		return problems
	},
}
