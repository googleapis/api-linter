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

package aip0134

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// Update methods should have an HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "http-body"),
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		fieldName := strcase.SnakeCase(m.GetName()[6:])
		// Establish that the RPC has HTTP body equal to fieldName.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Body != fieldName {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Update methods should have an HTTP body equal to `%q`.", fieldName),
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
