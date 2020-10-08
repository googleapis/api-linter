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

func TestAllowMissing(t *testing.T) {
	problems := testutils.Problems{{Message: "include `bool allow_missing`"}}
	for _, test := range []struct {
		name         string
		Style        string
		AllowMissing string
		problems     testutils.Problems
	}{
		{"IgnoredNotDF", "", "", nil},
		{"ValidIncluded", "style: DECLARATIVE_FRIENDLY", "bool allow_missing = 2;", nil},
		{"Invalid", "style: DECLARATIVE_FRIENDLY", "", problems},
		{"InvalidWrongType", "style: DECLARATIVE_FRIENDLY", "string allow_missing = 2;", problems},
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
						{{.Style}}
					};
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(allowMissing.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
