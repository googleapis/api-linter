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

package utils

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestDeclarativeFriendlyMessage(t *testing.T) {
	// Test the cases where a google.api.resource annotation is present.
	for _, test := range []struct {
		name  string
		Style string
		want  bool
	}{
		{"True", "style: DECLARATIVE_FRIENDLY", true},
		{"FalseNoStyle", "", false},
		{"FalseOtherStyle", "style: STYLE_UNSPECIFIED", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						{{.Style}}
					};
				}

				message CreateBookRequest {
					Book book = 1;
				}

				service Library {
					rpc CreateBook(CreateBookRequest) returns (Book);
				}
			`, test)
			for i := 0; i < f.Messages().Len(); i++ {
				m := f.Messages().Get(i)
				t.Run(string(m.Name()), func(t *testing.T) {
					if got := IsDeclarativeFriendlyMessage(m); got != test.want {
						t.Errorf("Got %v, expected %v.", got, test.want)
					}
				})
			}
		})
	}

	// Test the case where the google.api.resource annotation is not present.
	t.Run("NotResource", func(t *testing.T) {
		f := testutils.ParseProto3Tmpl(t, "message Book {}", nil)
		m := f.Messages().Get(0)
		if IsDeclarativeFriendlyMessage(m) {
			t.Errorf("Got true, expected false.")
		}
	})
}

func TestDeclarativeFriendlyMethod(t *testing.T) {
	// We need different templates for different situations.
	//
	// Note: The Book resource itself is always present and omitted here to
	// avoid excess repetition; it is appended to the templates in the body of
	// the test.
	tmpls := map[string]string{
		// The basic template just returns the resource with no frills.
		"basic": `
			service Library {
				rpc GetBook(GetBookRequest) returns (Book);
			}

			message GetBookRequest {}
		`,

		// The LRO template returns the resource, but as the result of an LRO
		// that has to be resolved first.
		"lro": `
			import "google/longrunning/operations.proto";

			service Library {
				rpc GetBook(GetBookRequest) returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "Book"
					};
				}
			}

			message GetBookRequest {}
		`,

		// The List template returns a normal list response.
		"list": `
			service Library {
				rpc ListBooks(ListBooksRequest) returns (ListBooksResponse);
			}

			message ListBooksRequest {}
			message ListBooksResponse {
				repeated Book books = 1;
			}
		`,

		// The Delete template returns a normal delete response.
		"delete": `
			import "google/protobuf/empty.proto";
			service Library {
				rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty);
			}
			message DeleteBookRequest {}
		`,

		// The custom method template is a straightforward custom method
		// with no direct reference to Book.
		"custom": `
			service Library {
				rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResponse);
			}

			message ArchiveBookRequest {}
			message ArchiveBookResponse {}
		`,
	}

	for key, tmpl := range tmpls {
		t.Run(key, func(t *testing.T) {
			for _, test := range []struct {
				name string
				want bool
			}{
				{"true", true},
				{"false", false},
			} {
				t.Run(test.name, func(t *testing.T) {
					// Set the style of the resource to DECLARATIVE_FRIENDLY if that
					// is the expected result.
					s := struct{ Style string }{Style: ""}
					if test.want == true {
						s.Style = "style: DECLARATIVE_FRIENDLY"
					}

					// Parse the template and test the method.
					f := testutils.ParseProto3Tmpl(t, fmt.Sprintf(`
						import "google/api/resource.proto";

						%s

						message Book {
							option (google.api.resource) = {
								type: "library.googleapis.com/Book"
								{{.Style}}
							};
						}
					`, tmpl), s)
					m := f.Services().Get(0).Methods().Get(0)
					if got := IsDeclarativeFriendlyMethod(m); got != test.want {
						t.Errorf("Got %v, expected %v.", got, test.want)
					}
				})
			}
		})
	}

	// Test an edge case where the LRO response is not found.
	t.Run("lro/not-found", func(t *testing.T) {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/longrunning/operations.proto";
			service Library {
				rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "Shrug"
					};
				}
			}
			message CreateBookRequest {}
		`, nil)
		m := f.Services().Get(0).Methods().Get(0)
		want := false
		if got := IsDeclarativeFriendlyMethod(m); got != want {
			t.Errorf("Got %v, expected %v.", got, want)
		}
	})
}

