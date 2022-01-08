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
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// Update methods should use the word "update", not synonyms.
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName(134, "synonyms"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return m.GetName() != "SetIamPolicy"
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		name := m.GetName()
		for _, syn := range []string{"Patch", "Put", "Set"} {
			if strings.HasPrefix(name, syn) {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						`%q can be a synonym for "Update". Should this be a Update method?`,
						syn,
					),
					Descriptor: m,
					Location:   locations.DescriptorName(m),
					Suggestion: strings.Replace(name, syn, "Update", 1),
				}}
			}
		}
		return nil
	},
}
