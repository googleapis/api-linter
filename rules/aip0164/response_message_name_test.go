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
			message {{.RespTypeName}} {}
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
		{"ValidResource", tmpl["sync"], "UndeleteBook", "Book", testutils.Problems{}},
		{"ValidLROResource", tmpl["lro"], "UndeleteBook", "Book", testutils.Problems{}},
		{"Invalid", tmpl["sync"], "UndeleteBook", "UndeleteBookResponse", testutils.Problems{{Suggestion: "Book"}}},
		{"InvalidLRO", tmpl["lro"], "UndeleteBook", "UndeleteBookResponse", testutils.Problems{{Suggestion: "Book"}}},
		{"Irrelevant", tmpl["sync"], "DestroyBook", "DestroyBookResponse", testutils.Problems{}},
		{"IrrelevantLRO", tmpl["lro"], "DestroyBook", "DestroyBookResponse", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-164 Undelete method
			file := testutils.ParseProto3Tmpl(t, test.tmpl, test)

			// Run the lint rule, and establish that it returns the expected problems.
			method := file.GetServices()[0].GetMethods()[0]
			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
