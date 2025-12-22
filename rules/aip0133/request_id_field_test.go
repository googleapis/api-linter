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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestRequestIDField(t *testing.T) {
	problems := testutils.Problems{{Message: "`string book_id`"}}
	for _, test := range []struct {
		name     string
		IDField  string
		problems testutils.Problems
	}{
		{"Valid", "string book_id = 2;", nil},
		{"InvalidMissing", "", problems},
		{"InvalidType", "bytes book_id = 2;", problems},
		{"InvalidRepeated", "repeated string book_id = 2;", problems},
		// request_id allowlist: skip check if valid request_id is present (AIP-133 exceptions)
		{"ValidRequestIdOnly", "string request_id = 2;", nil},
		{"ValidBothIds", "string book_id = 2; string request_id = 4;", nil},
		// invalid request_id types should NOT bypass the check
		{"InvalidRequestIdWrongType", "bytes request_id = 2;", problems},
		{"InvalidRequestIdRepeated", "repeated string request_id = 2;", problems},
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
					};
				}

				message CreateBookRequest {
					string parent = 1;
					{{.IDField}}
					Book book = 3;
				}
			`, test)
			m := f.Messages().ByName("CreateBookRequest")
			if diff := test.problems.SetDescriptor(m).Diff(requestIDField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
