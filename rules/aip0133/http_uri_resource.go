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
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// The resource name used in the Create method's URI should match the name used
// in the resource definition.
var httpURIResource = &lint.MethodRule{
	Name: lint.NewRuleName(133, "http-uri-resource"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return utils.IsCreateMethod(m) && len(utils.GetHTTPRules(m)) > 0
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		problems := []lint.Problem{}

		// Extract the suffix of the URI path as the collection identifier.
		uriParts := strings.Split(utils.GetHTTPRules(m)[0].URI, "/")
		collectionName := uriParts[len(uriParts)-1]

		// Ensure that a collection identifier is provided.
		if collectionName == "" {
			return []lint.Problem{{
				Message:    "The URI path does not end in a collection identifier.",
				Descriptor: m,
			}}
		}

		// Go through each pattern in the resource and make sure it contains the
		// collection identifier.
		collectionName += "/"
		resource := utils.GetResource(m.GetOutputType())
		for _, pattern := range resource.GetPattern() {
			if !strings.Contains(pattern, collectionName) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Resource pattern should contain the collection identifier %q.", collectionName),
					Descriptor: m.GetOutputType(),
				})
			}
		}

		return problems
	},
}
