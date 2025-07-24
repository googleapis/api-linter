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

package aip0136

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		Method     string
		Body       string
		problems   testutils.Problems
	}{
		{"ValidPostBody", "ArchiveBook", "post", "*", testutils.Problems{}},
		{"ValidPostBodyBook", "ArchiveBook", "post", "book", testutils.Problems{}},
		{"InvalidPostBody", "ArchiveBook", "post", "random", testutils.Problems{{Message: `body: "*"`}}},
		{"InvalidPostNoBody", "ArchiveBook", "post", "", testutils.Problems{{Message: `body: "*"`}}},
		{"ValidGetNoBody", "ReadBook", "get", "", testutils.Problems{}},
		{"InvalidGetBody", "ReadBook", "get", "*", testutils.Problems{{Message: "should not set"}}},
		{"InvalidGetOtherBody", "ReadBook", "get", "Book", testutils.Problems{{Message: "should not set"}}},
		{"ValidStdMethod", "CreateBook", "post", "book", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							{{.Method}}: "/v1:frob"
							body: "{{.Body}}"
						};
					}
				}
				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}
			`, test)
			method := file.Services()[0].Methods()[0]
			problems := httpBody.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestHttpBody_exceptions(t *testing.T) {
	tests := []struct {
		testName string
		BodyType string
		problems testutils.Problems
	}{
		{"ValidPostBody", "google.api.HttpBody", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/httpbody.proto";
				service Library {
					rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResponse) {
						option (google.api.http) = {
							post: "/v1:frob"
							body: "body"
						};
					}
				}
				message ArchiveBookRequest {
					{{.BodyType}} body = 1;
				}
				message ArchiveBookResponse {}
			`, test)
			method := file.Services()[0].Methods()[0]
			problems := httpBody.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
