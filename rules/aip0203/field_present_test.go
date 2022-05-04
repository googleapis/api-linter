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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestPresent(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "valid case with field_behavior present",
			comment:  "field_behavior-present",
			field:    "string secret = 1 [(google.api.field_behavior) = INPUT_ONLY];",
			problems: nil,
		},
		{
			name:    "invalid case with field_behavior missing",
			comment: "field_behavior-missing",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
				import "google/api/field_behavior.proto";

				message Book {
					// Secrets to be stored in the book
					// {{.Comment}}
					{{.Field}}
				}
				`
			file := testutils.ParseProto3Tmpl(t, template, struct {
				Comment string
				Field   string
			}{test.comment, test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := fieldPresent.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
