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

package aip0136

import (
	pluralize "github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var httpVariables = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "http-variables"),
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		p := pluralize.NewClient()
		for _, http := utils.GetHTTPRules(m) {
			vars := http.GetVariables()

			// If there is a "name" variable, the noun should be present
			// in the RPC name.
			if name, ok := vars["name"]; ok {
				// Determine the resource.
				name = strings.TrimSuffix(name, "/*")
				segs := strings.Split(name, "/")
				resource := p.Singular(segs[len(segs) - 1])

				// Does the RPC name end in the singular name of the resource?
				// If not, complain.

			}

			if parent, ok := vars["parent"]; ok {
				// Determine the resource.
				// Because this is `parent`, we expect the plural version this time.
				// Furthermore, it is **outside** the actual parent variable.
				segs := strings.Split(strings.Split(http.GetPlainURI(), ":")[0], "/")
				resource := segs[len(segs) - 1]

				// Does the RPC name end in the plural name of the resource?
				// If not, complain.
			}
		}
		return nil
	},
}
