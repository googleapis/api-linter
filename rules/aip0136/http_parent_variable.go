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
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var httpParentVariable = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "http-parent-variable"),
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		p := pluralize.NewClient()
		for _, http := range utils.GetHTTPRules(m) {
			vars := http.GetVariables()

			// If there is a "parent" variable, the noun should be present
			// in the RPC name.
			if _, ok := vars["parent"]; ok {
				// Determine the resource.
				segs := strings.Split(strings.Split(http.GetPlainURI(), ":")[0], "/")
				plural := strcase.SnakeCase(segs[len(segs)-1])
				singular := p.Singular(plural)

				// Does the RPC name end in the singular or plural name of the resource?
				// If not, complain.
				snakeName := strcase.SnakeCase(m.GetName())
				if !strings.HasSuffix(snakeName, plural) && !strings.HasSuffix(snakeName, singular) {
					return []lint.Problem{{
						Message:    "The parent variable should only be used if the RPC noun matches the URI.",
						Descriptor: m,
						Location:   locations.MethodHTTPRule(m),
					}}
				}
			}
		}
		return nil
	},
}
