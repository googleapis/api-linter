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

func TestDeclarativeFriendlyRequired(t *testing.T) {
	t.Run("Resource", func(t *testing.T) {
		for _, test := range []struct {
			name     string
			Style    string
			Etag     string
			problems testutils.Problems
		}{
			{"ValidEtagNotDF", "", "string etag = 2;", nil},
			{"ValidNoEtagNotDF", "", "", nil},
			{"ValidEtagDF", "style: DECLARATIVE_FRIENDLY", "string etag = 2;", nil},
			{"InvalidNoEtagDF", "style: DECLARATIVE_FRIENDLY", "", testutils.Problems{{Message: "string etag"}}},
		} {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						{{.Style}}
					};

					string name = 1;
					{{.Etag}}
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(declarativeFriendlyRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		}
	})

	t.Run("Requests", func(t *testing.T) {
		for _, test := range []struct {
			name     string
			Style    string
			Etag     string
			problems testutils.Problems
		}{
			{"ValidEtagNotDF", "", "string etag = 2;", nil},
			{"ValidNoEtagNotDF", "", "", nil},
			{"ValidEtagDF", "style: DECLARATIVE_FRIENDLY", "string etag = 2;", nil},
			{"InvalidNoEtagDF", "style: DECLARATIVE_FRIENDLY", "", testutils.Problems{{Message: "string etag"}}},
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
						rpc CreateBook(CreateBookRequest) returns (Book) {
							option (google.api.http) = {
								post: "/v1/{parent=publishers/*}/books"
								body: "*"
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
							{{.Style}}
						};
						string name = 1;
						string etag = 2;
					}
					message GetBookRequest {
						string name = 1;
					}
					message CreateBookRequest {
						string parent = 1;
						Book book = 2;
						string book_id = 3;
					}
					message DeleteBookRequest {
						string name = 1;
						{{.Etag}}
					}
				`, test)
				m := f.FindMessage("DeleteBookRequest")
				if diff := test.problems.SetDescriptor(m).Diff(declarativeFriendlyRequired.Lint(f)); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})
}
