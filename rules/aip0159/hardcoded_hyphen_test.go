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

package aip0159

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHardcodedHyphen(t *testing.T) {
	for _, test := range []struct {
		name     string
		Parent   string
		problems testutils.Problems
	}{
		{"Valid", "publishers/*", testutils.Problems{}},
		{"Invalid", "publishers/-", testutils.Problems{{Message: "`-`"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/api/annotations.proto";
			service Library {
				rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
					option (google.api.http) = {
						get: "/v1/{parent={{.Parent}}}/books"
					};
				}
			}
			message ListBooksRequest {}
			message ListBooksResponse {}
		`, test)
		method := f.GetServices()[0].GetMethods()[0]
		if diff := test.problems.SetDescriptor(method).Diff(hardcodedHyphen.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	}
}
