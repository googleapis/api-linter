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

package aip0193

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"

	"github.com/googleapis/api-linter/rules/internal/utils"
)

const expectedType = "google.rpc.Status"

// Set of field names that are expected to be typed with the standard status type.
var statusFieldNames = map[string]struct{}{
	"error":    struct{}{},
	"errors":   struct{}{},
	"status":   struct{}{},
	"statuses": struct{}{},
	"warning":  struct{}{},
	"warnings": struct{}{},
}

var responseStatusTypeCheck = &lint.MethodRule{
	Name: lint.NewRuleName(193, "response-status-type"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var p []lint.Problem
		for _, f := range m.GetOutputType().GetFields() {
			_, isStatusField := statusFieldNames[f.GetName()]
			if isStatusField && utils.GetTypeName(f) != expectedType {
				p = append(p, lint.Problem{
					Message:    fmt.Sprintf("Services must return `%s` message when API error occur.", expectedType),
					Descriptor: f,
					Suggestion: expectedType,
					Location:   locations.FieldType(f),
				})
			}
		}
		return p
	},
}
