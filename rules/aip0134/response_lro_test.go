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

package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseLRO(t *testing.T) {
	problems := testutils.Problems{{Suggestion: "google.longrunning.Operation"}}
	for _, test := range []struct {
		name         string
		Style        string
		ResponseType string
		problems     testutils.Problems
	}{
		{"ValidNotDF", "", "Book", nil},
		{"ValidLRO", "style: DECLARATIVE_FRIENDLY", "google.longrunning.Operation", nil},
		{"InvalidDFSync", "style: DECLARATIVE_FRIENDLY", "Book", problems},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns ({{.ResponseType}});
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						{{.Style}}
					};
				}

				message UpdateBookRequest {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(responseLRO.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
