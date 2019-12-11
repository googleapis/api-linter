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

package aip0127

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceNameExtraction(t *testing.T) {
	for _, test := range []struct {
		name     string
		uri      string
		problems testutils.Problems
	}{
		{"Valid", "/v1/{name=publishers/*/books/*}", testutils.Problems{}},
		{"VersioningTool", "/{$api_version}/{name=publishers/*/books/*}", testutils.Problems{}},
		{"Invalid", "/v1/publishers/{publisher_id}/books/{book_id}", testutils.Problems{{Message: "full resource name"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3String(t, strings.ReplaceAll(`
				import "google/api/annotations.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "{{.URI}}"
						};
					}
				}
				message GetBookRequest {}
				message Book {}
			`, "{{.URI}}", test.uri))
			method := f.GetServices()[0].GetMethods()[0]
			problems := resourceNameExtraction.Lint(f)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
