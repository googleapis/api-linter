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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestIDField(t *testing.T) {
	problems := testutils.Problems{{Message: "`string book_id`"}}
	for _, test := range []struct {
		name     string
		Style    string
		IDField  string
		problems testutils.Problems
	}{
		{"ValidNotDF", "", "", nil},
		{"ValidClientSpecified", "style: DECLARATIVE_FRIENDLY", "string book_id = 3;", nil},
		{"InvalidDF", "style: DECLARATIVE_FRIENDLY", "", problems},
		{"InvalidType", "style: DECLARATIVE_FRIENDLY", "bytes book_id = 3;", problems},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				service Library {
					rpc CreateBook(CreateBookRequest) returns (Book);
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						{{.Style}}
					};
				}

				message CreateBookRequest {
					string parent = 1;
					Book book = 2;
					{{.IDField}}
				}
			`, test)
			m := f.FindMessage("CreateBookRequest")
			if diff := test.problems.SetDescriptor(m).Diff(requestIDField.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
