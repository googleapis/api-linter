// Copyright 2020 Google LLC
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

package aip0128

import (
	"fmt"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceAnnotationsField = &lint.MessageRule{
	Name:   lint.NewRuleName(128, "resource-annotations-field"),
	OnlyIf: isDeclarativeFriendlyResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		f := m.FindFieldByName("annotations")
		if f == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Declarative-friendly %q has no `annotations` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		if t := utils.GetTypeName(f); t != "map<string, string>" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The `annotations` field should be a `map<string, string>`, not %s.", t),
				Descriptor: f,
				Location:   locations.FieldType(f),
				Suggestion: "map<string, string>",
			}}
		}

		return nil
	},
}
