// Copyright 2019 Google LLC
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

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName      string
		ResourceField string
		Body          string
		MethodName    string
		problems      testutils.Problems
	}{
		{"Valid", "Book book = 1;", "book", "CreateBook", nil},
		{"Valid", "Book textbook = 1;", "textbook", "CreateBook", nil},
		{"Valid", "", "book", "CreateBook", nil}, // valid for http body rule check, but it will fail under resource fail rule check
		{"Invalid_BodyMissing", "Book book = 1;", "", "CreateBook", testutils.Problems{{Message: "Post methods should have an HTTP body"}}},
		{"Invalid_BodyMismatch", "Book book = 1;", "abook", "CreateBook", testutils.Problems{{Message: `The content of body "abook" must map to the resource field "book" in the request message`}}},
		{"Irrelevant", "Book book = 1;", "book", "CreateBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/books"
							body: "{{.Body}}"
						};
					}
				}
				message Book {}
				message {{.MethodName}}Request {
					{{.ResourceField}}
				}
			`, test)
			method := file.Services()[0].Methods()[0]
			problems := httpBody.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
