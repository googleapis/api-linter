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

func TestHasAnnotation(t *testing.T) {
	ann := `option (google.api.http).get = "/v1/foo/";`
	tests := []struct {
		name       string
		Istream    string
		Ostream    string
		annotation string
		problems   testutils.Problems
	}{
		{"ValidUU", "", "", ann, testutils.Problems{}},
		{"InvalidUU", "", "", "", testutils.Problems{{Message: "google.api.http"}}},
		{"ValidUS", "", "stream ", ann, testutils.Problems{}},
		{"InvalidUS", "", "stream ", "", testutils.Problems{{Message: "google.api.http"}}},
		{"ValidSU", "stream ", "", ann, testutils.Problems{}},
		{"InvalidSU", "stream ", "", "", testutils.Problems{{Message: "google.api.http"}}},
		{"ValidSS", "stream ", "stream ", "", testutils.Problems{}},
		{"InvalidSS", "stream ", "stream ", ann, testutils.Problems{{Message: "google.api.http"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use of strings.ReplaceAll here allows a replacement using quotes,
			// which Go templates has no way to get around.
			f := testutils.ParseProto3Tmpl(t, strings.ReplaceAll(`
				import "google/api/annotations.proto";
				service Library {
					rpc ReadBook({{.Istream}}ReadBookRequest) returns ({{.Ostream}}ReadBookResponse) {
						{{.Annotation}}
					}
				}
				message ReadBookRequest {}
				message ReadBookResponse {}
			`, "{{.Annotation}}", test.annotation), test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(hasAnnotation.Lint(m)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
