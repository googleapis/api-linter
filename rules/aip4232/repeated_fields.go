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
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Only the last field in a method signature can be repeated.
var repeatedFields = &lint.MethodRule{
	Name:   lint.NewRuleName(4232, "repeated-fields"),
	OnlyIf: hasMethodSignatures,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		var problems []lint.Problem
		sigs := utils.GetMethodSignatures(m)
		in := m.Input()

		for _, sig := range sigs {
			for _, name := range sig {
				if !strings.Contains(name, ".") {
					continue
				}
				chain := strings.Split(name, ".")

				// Exclude the last one from the look up since it can be repeated.
				for i := range chain[:len(chain)-1] {
					n := strings.Join(chain[:i+1], ".")
					if f := utils.FindFieldDotNotation(in, n); f != nil && f.IsRepeated() {
						problems = append(problems, lint.Problem{
							Message:    fmt.Sprintf("Field %q in method_signature %q: only the last field in a signature argument can be repeated.", name, strings.Join(sig, ",")),
							Descriptor: m,
						})
					}
				}
			}
		}

		return problems
	},
}
