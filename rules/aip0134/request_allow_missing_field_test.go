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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestAllowMissing(t *testing.T) {
	const singletonPattern = `books/{book}/settings`
	const nonSingletonPattern = `books/{book}`
	problems := testutils.Problems{{Message: "include a singular `bool allow_missing`"}}
	for _, test := range []struct {
		name         string
		Style        string
		Pattern      string
		AllowMissing string
		problems     testutils.Problems
	}{
		{"IgnoredNotDF", "", nonSingletonPattern, "", nil},
		{"ValidIncluded", "style: DECLARATIVE_FRIENDLY", nonSingletonPattern, "bool allow_missing = 2;", nil},
		{"Invalid", "style: DECLARATIVE_FRIENDLY", nonSingletonPattern, "", problems},
		{"InvalidWrongType", "style: DECLARATIVE_FRIENDLY", nonSingletonPattern, "string allow_missing = 2;", problems},
		{"InvalidRepeated", "style: DECLARATIVE_FRIENDLY", nonSingletonPattern, "repeated bool allow_missing = 2;", problems},
		{"IgnoredSingleton", "style: DECLARATIVE_FRIENDLY", singletonPattern, "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (Book);
				}

				message UpdateBookRequest {
					Book book = 1;
					{{.AllowMissing}}
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "{{.Pattern}}"
						{{.Style}}
					};
				}
			`, test)
			m := f.Messages().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(allowMissing.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
