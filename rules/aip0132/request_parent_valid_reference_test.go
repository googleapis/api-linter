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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestParentValidReference(t *testing.T) {
	for _, test := range []struct {
		name          string
		ReferenceType string
		problems      testutils.Problems
	}{
		{"Valid", "type: \"library.googleapis.com/Publisher\"", testutils.Problems{}},
		{"Invalid", "type: \"library.googleapis.com/Book\"", testutils.Problems{{Message: "reference the parent(s)"}}},
		{"IgnoreChildType", "child_type: \"library.googleapis.com/Book\"", testutils.Problems{}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			message ListBooksRequest {
				string parent = 1 [(google.api.resource_reference) = {
					{{.ReferenceType}}
				}];
			}

			message ListBooksResponse {
				repeated Book books = 1;
			}

			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};

				string name = 1;
			}
		`, test)
		field := f.Messages()[0].Fields()[0]
		if diff := test.problems.SetDescriptor(field).Diff(requestParentValidReference.Lint(f)); diff != "" {
			t.Error(diff)
		}
	}
}
