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

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

var resourceReconcilingField = &lint.MessageRule{
	Name:   lint.NewRuleName(128, "resource-reconciling-field"),
	OnlyIf: isDeclarativeFriendlyResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		f := m.FindFieldByName("reconciling")
		if f == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Declarative-friendly %q has no `reconciling` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		if f.GetType() != builder.FieldTypeBool().GetType() || f.IsRepeated() {
			return []lint.Problem{{
				Message:    "The `reconciling` field should be a singular bool.",
				Descriptor: f,
				Location:   locations.FieldType(f),
				Suggestion: "bool",
			}}
		}

		return nil
	},
}
