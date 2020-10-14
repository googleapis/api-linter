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
			import "google/api/resource.proto";
			import "google/protobuf/empty.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespTypeName}});
			}
			message {{.MethodName}}Request {}
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
					{{.Style}}
				};
			}
			{{ if (ne .RespTypeName "google.protobuf.Empty") }}{{ if (ne .RespTypeName "Book") }}
			message {{.RespTypeName}} {}
			{{ end }}{{ end }}
		`,
		"lro": `
			package test;
			import "google/api/resource.proto";
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
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
					{{.Style}}
				};
			}
		`,
	}

	// Prepare problems.
	problems := map[string]map[string]testutils.Problems{
		"book": {
			"sync": {{Suggestion: "Book", Message: "of the resource"}},
			"lro":  {{Message: "of the resource"}},
		},
		"empty": {
			"sync": {{Suggestion: "google.protobuf.Empty", Message: "of Empty or the resource"}},
			"lro":  {{Message: "of Empty or the resource"}},
		},
		"none": {"sync": nil, "lro": nil},
	}

	// Set up the testing permutations.
	tests := []struct {
		name         string
		MethodName   string
		RespTypeName string
		Style        string
		problems     map[string]testutils.Problems
	}{
		{"ValidEmpty", "DeleteBook", "google.protobuf.Empty", "", problems["none"]},
		{"ValidResource", "DeleteBook", "Book", "", problems["none"]},
		{"Invalid", "DeleteBook", "DeleteBookResponse", "", problems["empty"]},
		{"InvalidEmptyDF", "DeleteBook", "google.protobuf.Empty", "style: DECLARATIVE_FRIENDLY", problems["book"]},
		{"Irrelevant", "DestroyBook", "DestroyBookResponse", "", problems["none"]},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, tmplName := range []string{"sync", "lro"} {
				t.Run(tmplName, func(t *testing.T) {
					// Create a minimal service with a AIP-135 Delete method
					file := testutils.ParseProto3Tmpl(t, tmpl[tmplName], test)

					// Run the lint rule, and establish that it returns the expected problems.
					method := file.GetServices()[0].GetMethods()[0]
					problems := responseMessageName.Lint(file)
					if diff := test.problems[tmplName].SetDescriptor(method).Diff(problems); diff != "" {
						t.Errorf(diff)
					}
				})
			}
		})
	}
}
