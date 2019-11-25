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

package aip0231

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Batch Get methods should use the HTTP GET verb.
var httpVerb = &descrule.MethodRule{
	RuleName:   lint.NewRuleName(231, "http-method"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP GET.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "GET" {
				return []lint.Problem{{
					Message:    "Batch Get methods must use the HTTP GET verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
