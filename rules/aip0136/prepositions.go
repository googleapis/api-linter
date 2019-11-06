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
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/data"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var noPrepositions = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0136", "prepositions"),
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		for _, word := range strings.Split(strcase.SnakeCase(m.GetName()), "_") {
			if data.Prepositions.Contains(word) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Method names should not include prepositions (%q).", word),
					Descriptor: m,
					Location:   locations.DescriptorName(m),
				})
			}
		}
		return
	},
}
