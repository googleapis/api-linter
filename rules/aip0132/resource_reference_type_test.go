// Copyright 2022 Google LLC
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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceReferenceType(t *testing.T) {
	bookResource := `
option (google.api.resource) = {
	type: "library.googleapis.com/Book"
	pattern: "shelves/{shelf}/books/{book}"
};
`

	// Set up testing permutations.
	tests := []struct {
		testName           string
		TypeName           string
		RefType            string
		ResourceAnnotation string
		problems           testutils.Problems
	}{
		{"ValidChildType", "library.googleapis.com/Book", "child_type", bookResource, nil},
		{"ValidType", "library.googleapis.com/Shelf", "type", bookResource, nil},
		{"InvalidType", "library.googleapis.com/Book", "type", bookResource, testutils.Problems{{Message: "not a `type`"}}},
		{"InvalidChildType", "library.googleapis.com/Shelf", "child_type", bookResource, testutils.Problems{{Message: "`child_type`"}}},
		{"SkipNonResource", "library.googleapis.com/Book", "child_type", "", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				service Library {
					rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {}
				}
				message ListBooksRequest {
					string parent = 1 [(google.api.resource_reference).{{ .RefType }} = "{{ .TypeName }}"];
				}
				message ListBooksResponse {
					repeated string unreachable = 2;
					repeated Book books = 1;
				}
				message Book {
					{{ .ResourceAnnotation }}
					string name = 1;
				}
			`, test)
			field := file.GetServices()[0].GetMethods()[0].GetInputType().FindFieldByName("parent")
			problems := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
