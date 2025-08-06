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
		{"RequestId", "GetBookRequest", "request_id", "string", testutils.Problems{}},
		{"ReadMask", "GetBookRequest", "read_mask", "google.protobuf.FieldMask", testutils.Problems{}},
		{"View", "GetBookRequest", "view", "View", testutils.Problems{}},
		{"Invalid", "GetBookRequest", "application_id", "string", testutils.Problems{{
			Message: "Unexpected field",
		}}},
		{"Irrelevant", "AcquireBookRequest", "application_id", "string", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				enum View {
					VIEW_UNSPECIFIED = 0;
					BASIC = 1;
					FULL = 2;
				}
				message {{.MessageName}} {
					string name = 1;
					{{.FieldType}} {{.FieldName}} = 2;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(1)

			// Run the lint rule, and establish that it returns the correct problems.
			wantProblems := test.problems.SetDescriptor(field)
			gotProblems := unknownFields.Lint(f)
			if diff := wantProblems.Diff(gotProblems); diff != "" {
				t.Error(diff)
			}
		})
	}
}