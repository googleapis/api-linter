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

package aip0131

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Get standard method should only have expected fields.
var standardFields = &lint.MessageRule{
	Name:   lint.NewRuleName(131, "request-name-field"),
	OnlyIf: isGetRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a name field is present.
		name := m.FindFieldByName("name")
		if name == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Method %q has no `name` field", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Ensure that the name field is the correct type.
		if name.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`name` field on Get RPCs should be a string",
				Descriptor: name,
			}}
		}

		return nil
	},
}
