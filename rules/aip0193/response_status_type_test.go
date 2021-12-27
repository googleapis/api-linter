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
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

var statusTypeProblemBase = lint.Problem{
	Message:    "Services must return `google.rpc.Status` message when API error occur.",
	Suggestion: "google.rpc.Status",
}

func TestResponseStatusType(t *testing.T) {
	for _, test := range []struct {
		name                   string
		NonServiceMessageField string
		RequestField           string
		ResponseField1         string
		ResponseField2         string
		expectedProblems       testutils.Problems
	}{
		{
			name:             "NoFields",
			expectedProblems: nil,
		},
		{
			name:                   "InvalidStatusTypeInIrrelevantMessage",
			NonServiceMessageField: "int32 error = 1;",
			expectedProblems:       nil,
		},
		{
			name:             "NoErrorFields",
			RequestField:     "string query = 1;",
			ResponseField1:   "int32 results_per_page = 1;",
			expectedProblems: nil,
		},
		{
			name:             "ErrorFieldInRequest",
			RequestField:     "int32 error = 1;",
			expectedProblems: nil,
		},
		{
			name:             "ValidStatusTypeInResponse",
			ResponseField1:   "google.rpc.Status error = 1;",
			expectedProblems: nil,
		},
		{
			name:             "InvalidStatusTypeInResponse",
			ResponseField1:   "int32 error = 1;",
			expectedProblems: testutils.Problems{statusTypeProblemBase},
		},
		{
			name:             "2InvalidStatusTypesInResponse",
			ResponseField1:   "int32 error = 1;",
			ResponseField2:   "string status = 2;",
			expectedProblems: testutils.Problems{statusTypeProblemBase, statusTypeProblemBase},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			        import "google/rpc/status.proto";

				message Store {
					{{.NonServiceMessageField}}
				}

				message GetBookRequest{
					{{.RequestField}}
				}

				message GetBookResponse {
					{{.ResponseField1}}
					{{.ResponseField2}}
				}

				service BookService {
					rpc GetBook(GetBookRequest) returns (GetBookResponse) { }
				}
			`, test)

			// Get field descriptors of problematic fields.
			var problemFields []*desc.FieldDescriptor
			for i := range test.expectedProblems {
				field := f.GetServices()[0].GetMethods()[0].GetOutputType().GetFields()[i]
				test.expectedProblems[i].Descriptor = field
				problemFields = append(problemFields, field)
			}

			p := responseStatusTypeCheck.Lint(f)

			// Test the number of produced problems.
			if got, want := len(p), len(test.expectedProblems); got != want {
				t.Fatalf("Incorrect number of problems, got %d, want %d", got, want)
			}

			// Basic check of problem contents.
			diff := test.expectedProblems.Diff(p)
			if diff != "" {
				t.Errorf(diff)
			}

			// Check locations.
			for i := range test.expectedProblems {
				want := locations.FieldType(problemFields[i])
				if got := p[i].Location; got != want {
					t.Errorf("Incorrect location, got %v, want %v", got, want)
				}
			}
		})
	}
}
