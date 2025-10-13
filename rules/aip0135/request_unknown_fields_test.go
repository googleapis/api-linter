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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		MessageName string
		FieldName   string
		FieldType   string
		problems    testutils.Problems
	}{
		{"Force", "DeleteBookRequest", "force", "bool", testutils.Problems{}},
		{"Etag", "DeleteBookRequest", "etag", "string", testutils.Problems{}},
		{"AllowMissing", "DeleteBookRequest", "allow_missing", "bool", testutils.Problems{}},
		{"RequestId", "DeleteBookRequest", "request_id", "string", testutils.Problems{}},
		{"ValidateOnly", "DeleteBookRequest", "validate_only", "bool", testutils.Problems{}},
		{"ValidView", "DeleteBookRequest", "view", "BookView", testutils.Problems{}},
		{"SuffixedView", "DeleteBookRequest", "custom_view", "BookView", testutils.Problems{}},
		{"Invalid", "DeleteBookRequest", "application_id", "string", testutils.Problems{{
			Message: "Unexpected field",
		}}}, {"Irrelevant", "RemoveBookRequest", "application_id", "string", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				enum BookView {
					BOOK_VIEW_UNSPECIFIED = 0;
					BOOK_VIEW_BASIC = 1;
					BOOK_VIEW_FULL = 2;
				}
				message {{.MessageName}} {
					string name = 1;
					{{.FieldType}} {{.FieldName}} = 2;
				}
			`, test)

			// Run the lint rule, and establish that it returns the correct problems.
			problems := unknownFields.Lint(f)
			if diff := test.problems.SetDescriptor(f.Messages().Get(0).Fields().Get(1)).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
