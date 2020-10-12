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

package aip0163

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestDeclarativeFriendlyRequired(t *testing.T) {
	problems := testutils.Problems{{Message: "`bool validate_only`"}}
	for _, test := range []struct {
		name         string
		Style        string
		ValidateOnly string
		problems     testutils.Problems
	}{
		{"ValidNotDF", "", "", nil},
		{"ValidDF", "style: DECLARATIVE_FRIENDLY", "bool validate_only = 2;", nil},
		{"InvalidDFWrongType", "style: DECLARATIVE_FRIENDLY", "int32 validate_only = 2;", problems},
		{"InvalidDFMissing", "style: DECLARATIVE_FRIENDLY", "", problems},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";

				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "/v1/{name=publishers/*/books/*}"
						};
					}

					rpc DeleteBook(DeleteBookRequest) returns (Book) {
						option (google.api.http) = {
							delete: "/v1/{name=publishers/*/books/*}"
						};
					}
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						{{.Style}}
					};
					string name = 1;
				}

				message GetBookRequest {
					string name = 1;
				}

				message DeleteBookRequest {
					string name = 1;
					{{.ValidateOnly}}
				}
			`, test)
			dbr := f.FindMessage("DeleteBookRequest")
			if diff := test.problems.SetDescriptor(dbr).Diff(declarativeFriendlyRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
