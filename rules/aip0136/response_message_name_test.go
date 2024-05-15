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
	// Set up the testing permutations.
	tests := []struct {
		testName            string
		MethodName          string
		RespMessageName     string
		OperatingOnResource bool
		problems            testutils.Problems
	}{
		{"Valid Resource", "ArchiveBook", "Book", true, testutils.Problems{}},
		{"Invalid Resource", "ArchiveBook", "Author", true, testutils.Problems{{Suggestion: "ArchiveBookResponse or Book"}}},
		{"Valid Response Suffix Stateful", "ArchiveBook", "ArchiveBookResponse", true, testutils.Problems{}},
		{"Valid Response Suffix Stateless", "ArchiveBook", "ArchiveBookResponse", false, testutils.Problems{}},
		{"Invalid Response Suffix", "ArchiveBook", "ArchiveBookResp", true, testutils.Problems{{Suggestion: "ArchiveBookResponse or Book"}}},
		{"Unable To Find Resource", "ArchiveBook", "ArchiveBookResp", false, testutils.Problems{{Suggestion: "ArchiveBookResponse"}}},
		{"Irrelevant", "DeleteBook", "DeleteBookResp", true, testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			package test;
			import "google/api/resource.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespMessageName}});
			}
			message {{.MethodName}}Request {
				{{ if (.OperatingOnResource) }}
				// The book to archive.
				// Format: publishers/{publisher}/books/{book}
				string name = 1 [
				  (google.api.resource_reference) = {
					type: "library.googleapis.com/Book"
				  }];
				{{ end }}
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
			{{ if and (ne .RespMessageName "Book") (ne .RespMessageName "Author") }}
			message {{.RespMessageName}} {}
			{{ end }}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
