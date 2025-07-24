// Copyright 2022 Google LLC
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

func TestResourceReferenceType(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName string
		TypeName string
		RefType  string
		problems testutils.Problems
	}{
		{"ValidChildType", "library.googleapis.com/Book", "child_type", nil},
		{"ValidType", "library.googleapis.com/Shelf", "type", nil},
		{"InvalidType", "library.googleapis.com/Book", "type", testutils.Problems{{Message: "not a `type`"}}},
		{"InvalidChildType", "library.googleapis.com/Shelf", "child_type", testutils.Problems{{Message: "`child_type`"}}},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";
				service Library {
					rpc CreateBook(CreateBookRequest) returns (Book) {}
				}
				message CreateBookRequest {
					string parent = 1 [(google.api.resource_reference).{{ .RefType }} = "{{ .TypeName }}"];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "shelves/{shelf}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			field := file.Services()[0].Methods()[0].Input().FindFieldByName("parent")
			problems := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestResourceReferenceTypeLRO(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName     string
		TypeName     string
		RefType      string
		ResponseType string
		problems     testutils.Problems
	}{
		{"ValidChildType", "library.googleapis.com/Book", "child_type", "Book", nil},
		{"ValidChildTypeLRO", "library.googleapis.com/Book", "child_type", "Book", nil},
		{"ValidType", "library.googleapis.com/Shelf", "type", "Book", nil},
		{"InvalidType", "library.googleapis.com/Book", "type", "Book", testutils.Problems{{Message: "not a `type`"}}},
		{"InvalidChildType", "library.googleapis.com/Shelf", "child_type", "Book", testutils.Problems{{Message: "`child_type`"}}},
		{"SkipInvalidUnresolvableResponseType", "library.googleapis.com/Shelf", "child_type", "Foo", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";
				service Library {
					rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{ .ResponseType }}"
							metadata_type: "Book"
						};
					}
				}
				message CreateBookRequest {
					string parent = 1 [(google.api.resource_reference).{{ .RefType }} = "{{ .TypeName }}"];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "shelves/{shelf}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			field := file.Services()[0].Methods()[0].Input().FindFieldByName("parent")
			problems := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
