// Copyright 2021 Google LLC
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

package aip4232

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Only the last field in a method signature can be repeated.
var repeatedFields = &lint.MethodRule{
	Name:   lint.NewRuleName(4232, "repeated-fields"),
	OnlyIf: hasMethodSignatures,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var problems []lint.Problem
		sigs := utils.GetMethodSignatures(m)
		in := m.GetInputType()

		for _, sig := range sigs {
			for i, name := range sig {
				// Skip the last field in the signature, it can be repeated.
				if i == len(sig)-1 {
					break
				}

				if f := in.FindFieldByName(name); f != nil && f.IsRepeated() {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Field %q in method_signature %q: only the last field can be repeated.", name, strings.Join(sig, ",")),
						Descriptor: m,
					})
				}
			}
		}

		return problems
	},
}
