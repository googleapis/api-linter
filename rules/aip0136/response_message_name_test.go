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

package aip0136

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseMessageName(t *testing.T) {
	t.Run("Response Suffix", func(t *testing.T) {
		// Set up the testing permutations.
		tests := []struct {
			testName        string
			MethodName      string
			RespMessageName string
			problems        testutils.Problems
		}{
			{"Valid", "ArchiveBook", "ArchiveBookResponse", testutils.Problems{}},
			{"Invalid", "ArchiveBook", "ArchiveBookResp", testutils.Problems{{Message: "not \"ArchiveBookResp\"."}}},
		}

		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				file := testutils.ParseProto3Tmpl(t, `
				package test;
				import "google/api/resource.proto";

				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespMessageName}});
				}

				message {{.MethodName}}Request {}
				message {{.RespMessageName}} {}
				`, test)
				method := file.GetServices()[0].GetMethods()[0]
				problems := responseMessageName.Lint(file)
				if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
					t.Error(diff)
				}
			})
		}
	})

	t.Run("Response Suffix - LRO", func(t *testing.T) {
		// Set up the testing permutations.
		tests := []struct {
			testName    string
			MethodName  string
			MessageName string
			problems    testutils.Problems
		}{
			{"Valid", "ArchiveBook", "ArchiveBookResponse", testutils.Problems{}},
			{"Invalid", "ArchiveBook", "ArchiveBookResp", testutils.Problems{{Message: "not \"ArchiveBookResp\"."}}},
		}

		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				file := testutils.ParseProto3Tmpl(t, `
				package test;
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";

				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.MessageName}}"
							metadata_type: "OperationMetadata"
						};
					};
				}
				message {{.MethodName}}Request {}
				message {{.MessageName}} {}
				message OperationMetadata {}
				`, test)
				method := file.GetServices()[0].GetMethods()[0]
				problems := responseMessageName.Lint(file)
				if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
					t.Error(diff)
				}
			})
		}
	})

	t.Run("Resource", func(t *testing.T) {
		// Set up the testing permutations.
		tests := []struct {
			testName        string
			MethodName      string
			RespMessageName string
			LRO             bool
			problems        testutils.Problems
		}{
			{"Valid", "ArchiveBook", "Book", false, testutils.Problems{}},
			{"Valid LRO", "ArchiveBook", "Book", true, testutils.Problems{}},
			{"Invalid", "ArchiveBook", "Author", false, testutils.Problems{{Message: "not \"Author\"."}}},
			{"Invalid LRO", "ArchiveBook", "Author", true, testutils.Problems{{Message: "not \"Author\"."}}},
		}

		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				file := testutils.ParseProto3Tmpl(t, `
				package test;

				import "google/api/annotations.proto";
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";

				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{ if .LRO }}google.longrunning.Operation{{ else }}{{.RespMessageName}}{{ end }}) {
						option (google.api.http) = {
							post: "/v1/{name=publishers/*/books/*}:foo"
							body: "*"
						};
						{{ if .LRO }}
						option (google.longrunning.operation_info) = {
							response_type: "{{ .RespMessageName }}"
							metadata_type: "{{ .RespMessageName }}Metadata"
						};
						{{ end }}
					};
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
				}

				message Author {
					option (google.api.resource) = {
						type: "library.googleapis.com/Author"
						pattern: "authors/{author}"
					};
				}

				message {{.MethodName}}Request {
					// The book to operate on.
					// Format: publishers/{publisher}/books/{book}
					string name = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
				}
				`, test)
				method := file.GetServices()[0].GetMethods()[0]
				problems := responseMessageName.Lint(file)
				if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
					t.Error(diff)
				}
			})
		}
	})
}
