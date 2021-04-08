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

package aip0135

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestForceField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		FieldType   string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "DeleteBookRequest", "bool", "force", nil},
		{"Invalid", "DeleteBookRequest", "int32", "force", testutils.Problems{{Suggestion: "bool"}}},
		{"InvalidRepeated", "DeleteBookRequest", "repeated bool", "force", testutils.Problems{{Suggestion: "bool"}}},
		{"IrrelevantMessage", "RemoveBookRequest", "int32", "force", nil},
		{"IrrelevantField", "DeleteBookRequest", "string", "name", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			problems := requestForceField.Lint(f)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
