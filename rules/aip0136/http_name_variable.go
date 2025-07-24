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

	pluralize "github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

var httpNameVariable = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "http-name-variable"),
	OnlyIf: utils.IsCustomMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		p := pluralize.NewClient()
		for _, http := range utils.GetHTTPRules(m) {
			vars := http.GetVariables()

			// Special case: AIP-162 describes "revision" methods; the `name`
			// variable is appropriate (and mandated) for those.
			if strings.HasSuffix(m.Name(), "Revision") || strings.HasSuffix(m.Name(), "Revisions") {
				return nil
			}

			// If there is a "name" variable, the noun should be present
			// in the RPC name.
			if name, ok := vars["name"]; ok {
				// Determine the resource.
				name = strings.TrimSuffix(name, "/*")
				segs := strings.Split(name, "/")
				resource := strcase.SnakeCase(p.Singular(segs[len(segs)-1]))

				// Does the RPC name end in the singular name of the resource?
				// If not, complain.
				if !strings.HasSuffix(strcase.SnakeCase(m.Name()), resource) {
					return []lint.Problem{{
						Message:    "The name variable should only be used if the RPC noun matches the URI.",
						Descriptor: m,
						Location:   locations.MethodHTTPRule(m),
					}}
				}
			}
		}
		return nil
	},
}
