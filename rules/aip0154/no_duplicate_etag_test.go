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

package aip0154

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNoDuplicateEtag(t *testing.T) {
	for _, test := range []struct {
		name     string
		Etag     string
		problems testutils.Problems
	}{
		{"InvalidUpdate", "string etag = 2;", testutils.Problems{{Message: "omit etag", Suggestion: ""}}},
		{"ValidUpdate", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (Book);
					rpc DeleteBook(DeleteBookRequest) returns (Book);
				}

				message UpdateBookRequest {
					Book book = 1;
					{{.Etag}}
				}

				message DeleteBookRequest {
					string name = 1;
					string etag = 2;
				}

				message Book {
					string name = 1;
					string etag = 2;
				}
			`, test)
			m := f.Messages()
			if diff := test.problems.SetDescriptor(m.Get(0).Fields().ByName("etag")).Diff(noDuplicateEtag.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
