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

package aip0156

import (
	"fmt"
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var forbiddenMethods = &lint.MethodRule{
	Name: lint.NewRuleName(156, "forbidden-methods"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		// If the `name` variable in the URI ends in something other than
		// "*", that indicates that this is a singleton.
		//
		// For example:
		//   publishers/*/books/* -- not a singleton, many books
		//   publishers/*/settings -- a singleton; one settings object per publisher
		for _, http := range utils.GetHTTPRules(m) {
			if name, ok := http.GetVariables()["resource_name"]; ok {
				if !strings.HasSuffix(name, "*") {
					return true
				}
			}
		}
		return false
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Singletons should not use Create, List, or Delete.
		for _, badPrefix := range []string{"Create", "List", "Delete"} {
			if strings.HasPrefix(m.GetName(), badPrefix) {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Singletons must not define %q methods.", badPrefix),
					Descriptor: m,
				}}
			}
		}
		return nil
	},
}
