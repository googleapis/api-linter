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

package aip0164

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestRequestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		TestName    string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Etag", "UndeleteBookRequest", "etag", testutils.Problems{}},
		{"RequestId", "UndeleteBookRequest", "request_id", testutils.Problems{}},
		{"ValidateOnly", "UndeleteBookRequest", "validate_only", testutils.Problems{}},
		{"Invalid", "UndeleteBookRequest", "application_id", testutils.Problems{{
			Message: "Unexpected field",
		}}},
		{"Irrelevant", "RemoveBookRequest", "application_id", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					string name = 1;
					string {{.FieldName}} = 2;
				}
			`, test)
			message := f.Messages().Get(0)

			// Run the lint rule, and establish that it returns the correct problems.
			wantProblems := test.problems.SetDescriptor(message.Fields().Get(1))
			gotProblems := requestUnknownFields.Lint(f)
			if diff := wantProblems.Diff(gotProblems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
