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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestParentField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		name        string
		MessageName string
		FieldName   string
		FieldType   string
		problems    testutils.Problems
	}{
		{"Valid", "CreateBookRequest", "parent", "string", nil},
		{"InvalidType", "CreateBookRequest", "parent", "bytes", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantMessage", "AddBookRequest", "parent", "bytes", nil},
		{"IrrelevantField", "CreateBookRequest", "id", "bytes", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)

			problems := requestParentField.Lint(f)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
