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

func TestResponseMessageName(t *testing.T) {
	tmpl := map[string]string{
		"sync": `
			package test;
			import "google/protobuf/empty.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespTypeName}});
			}
			message {{.MethodName}}Request {}
			{{ if (ne .RespTypeName "google.protobuf.Empty") }}message {{.RespTypeName}} {}{{ end }}
		`,
		"lro": `
			package test;
			import "google/longrunning/operations.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request)
				    returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "{{.RespTypeName}}"
						metadata_type: "{{.MethodName}}Metadata"
					};
				}
			}
			message {{.MethodName}}Request {}
		`,
	}

	// Set up the testing permutations.
	tests := []struct {
		testName     string
		tmpl         string
		MethodName   string
		RespTypeName string
		problems     testutils.Problems
	}{
		{"ValidEmpty", tmpl["sync"], "DeleteBook", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidResource", tmpl["sync"], "DeleteBook", "Book", testutils.Problems{}},
		{"ValidLROEmpty", tmpl["lro"], "DeleteBook", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidLROResource", tmpl["lro"], "DeleteBook", "Book", testutils.Problems{}},
		{"Invalid", tmpl["sync"], "DeleteBook", "DeleteBookResponse", testutils.Problems{{Suggestion: "google.protobuf.Empty"}}},
		{"InvalidLRO", tmpl["lro"], "DeleteBook", "DeleteBookResponse", testutils.Problems{{Suggestion: "google.protobuf.Empty"}}},
		{"Irrelevant", tmpl["sync"], "DestroyBook", "DestroyBookResponse", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-135 Delete method
			file := testutils.ParseProto3Tmpl(t, test.tmpl, test)

			// Run the lint rule, and establish that it returns the expected problems.
			method := file.GetServices()[0].GetMethods()[0]
			problems := responseMessageName.Lint(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
