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

package aip0159

import (
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var hardcodedHyphen = &lint.MethodRule{
	Name: lint.NewRuleName(159, "hardcoded-hyphen"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		for _, http := range utils.GetHTTPRules(m) {
			if strings.Contains(http.GetPlainURI(), "/-/") {
				return []lint.Problem{{
					Message:    "URIs must not hard-code a `-` segment.",
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}
		return nil
	},
}
