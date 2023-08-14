// Copyright 2023 Google LLC
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

package aip0121

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/testutils"
)

// TestResourceMustSupportList tests the resourceMustSupportList
// lint rule by declaring a service proto, then declaring a
// google.api.resource message, then declaring non-List
// methods.
func TestResourceMustSupportList(t *testing.T) {
	for _, test := range []struct {
		name     string
		RPCs     string
		problems testutils.Problems
	}{
		{"ValidCreateList", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
			rpc CreateBook(CreateBookRequest) returns (Book) {};
		`, nil},
		{"ValidCreateListLRO", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
			rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "Book"
				};
			};
		`, nil},
		{"ValidListGet", `
			rpc GetBook(GetBookRequest) returns (Book) {};
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
		`, nil},
		{"ValidUpdateList", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
			rpc UpdateBook(UpdateBookRequest) returns (Book) {};
		`, nil},
		{"InvalidCreateOnly", `
			rpc CreateBook(CreateBookRequest) returns (Book) {};
		`, []lint.Problem{
			{Message: `resource "library.googleapis.com/Book"`},
		}},
		{"InvalidCreateOnlyLRO", `
			rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "Book"
				};
			};
		`, []lint.Problem{
			{Message: `resource "library.googleapis.com/Book"`},
		}},
		{"InvalidUpdateOnly", `
			rpc UpdateBook(UpdateBookRequest) returns (Book) {};
		`, []lint.Problem{
			{Message: `resource "library.googleapis.com/Book"`},
		}},
		{"InvalidGetOnly", `
			rpc GetBook(GetBookRequest) returns (Book) {};
		`, []lint.Problem{
			{Message: `resource "library.googleapis.com/Book"`},
		}},
		{"ValidIgnoreSingleton", `
			rpc GetBookCover(GetBookCoverRequest) returns (BookCover) {};
		`, nil},
		{"ValidIgnoreNonResource", `
			rpc GetBookCover(GetBookCoverRequest) returns (Other) {};
		`, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message BookCover {
					option (google.api.resource) = {
						type: "library.googleapis.com/BookCover"
						pattern: "books/{book}/bookCover"
						singular: "bookCover"
					};
				}

				message GetBookCoverRequest {
					string name = 1;
				}

				message CreateBookRequest {
					// The parent resource where this book will be created.
					// Format: publishers/{publisher}
					string parent = 1;

					// The book to create.
					Book book = 2;
				}

				message GetBookRequest {
					string name = 1;
				}

				message UpdateBookRequest {
					Book book = 1;
					google.protobuf.FieldMask update_mask = 2;
				}

				 message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
				 }

				 message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
				 }

				 message Other {
					string other = 1;
				 }
			`, test)
			s := file.GetServices()[0]
			got := resourceMustSupportList.Lint(file)
			if diff := test.problems.SetDescriptor(s).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
