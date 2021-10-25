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

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// All fields annotated as REQUIRED must be in every method_signature.
var requiredFields = &lint.MethodRule{
	Name:   lint.NewRuleName(4232, "required-fields"),
	OnlyIf: hasMethodSignatures,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var problems []lint.Problem
		sigs := utils.GetMethodSignatures(m)
		in := m.GetInputType()

		requiredFields := []string{}
		for _, f := range in.GetFields() {
			if utils.GetFieldBehavior(f).Contains("REQUIRED") {
				requiredFields = append(requiredFields, f.GetName())
			}
		}
		for i, sig := range sigs {
			sigset := stringset.New(sig...)
			if !sigset.Contains(requiredFields...) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Method signature %q missing at least one of the required fields: %q", strings.Join(sig, ","), strings.Join(requiredFields, ",")),
					Descriptor: m,
					Location:   locations.MethodSignature(m, i),
				})
			}
		}

		return problems
	},
}
