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

func TestOutputMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName     string
		MethodName   string
		RespTypeName string
		LRO          bool
		problems     testutils.Problems
	}{
		{"ValidResource", "CreateBook", "Book", false, testutils.Problems{}},
		{"ValidLRO", "CreateBook", "Book", true, testutils.Problems{}},
		{"Invalid", "CreateBook", "CreateBookResponse", false, testutils.Problems{{Suggestion: "Book"}}},
		{"InvalidLRO", "CreateBook", "CreateBookResponse", true, testutils.Problems{{Suggestion: "Book"}}},
		{"Irrelevant", "BuildBook", "BuildBookResponse", false, testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-134 Update method
			file := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns ({{ if .LRO }}google.longrunning.Operation{{ else }}{{ .RespTypeName }}{{ end }}) {
						{{ if .LRO -}}
						option (google.longrunning.operation_info) = {
							response_type: "{{.RespTypeName}}"
							metadata_type: "{{.MethodName}}Metadata"
						};
						{{ end -}}
					}
				}
				message {{.MethodName}}Request {}
				message {{.RespTypeName}} {}
			`, test)

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			problems := outputName.Lint(file)
			method := file.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
