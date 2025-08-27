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
	"strings"
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestHTTPURIResource(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "/v1/{parent=publishers/*}/books", "publishers/{publisher}/books/{book}", nil},
		{"ValidCustomLookalike", "/v1/{parent=publishers/*}/books:createAndCheckout", "publishers/{publisher}/books/{book}", nil},
		{"MethodMissingURIPath", "", "publishers/{publisher}/books/{book}", nil},
		{"MethodMissingCollectionURISuffix", "/v1/", "publishers/{publisher}/books/{book}", testutils.Problems{{Message: "The URI path does not end in a collection identifier."}}},
		{"ResourceMissingCollectionInPattern", "/v1/{parent=publishers/*}/books", "publishers/{publisher}", testutils.Problems{{Message: "Resource pattern should contain the collection identifier \"books/\"."}}},
		{"ResourceMissingCollectionCustomLookalike", "/v1/{parent=publishers/*}/books:createAndCheckout", "publishers/{publisher}", testutils.Problems{{Message: "Resource pattern should contain the collection identifier \"books/\"."}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";
				service Library {
					rpc CreateBook(CreateBookRequest) returns (Book) {
						option (google.api.http) = {
							post: "{{.URI}}"
						};
					}
				}
				message CreateBookRequest {}
				message Book {
					option (google.api.resource) = {
						pattern: "{{.Pattern}}"
					};
				}
			`, test)

			method := f.Services().Get(0).Methods().Get(0)
			var d protoreflect.Descriptor = method
			if strings.HasPrefix(test.testName, "Resource") {
				d = method.Output()
			}
			if diff := test.problems.SetDescriptor(d).Diff(httpURIResource.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
