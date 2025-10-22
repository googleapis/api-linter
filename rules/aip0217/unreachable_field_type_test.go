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

package aip0217

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestUnreachableFieldType(t *testing.T) {
	for _, test := range []struct {
		name     string
		Prefix   string
		problems testutils.Problems
	}{
		{"Valid", "repeated string", testutils.Problems{}},
		{"InvalidSingle", "string", testutils.Problems{{Message: "repeated"}}},
		{"InvalidType", "repeated ErrorMetadata", testutils.Problems{{Suggestion: "string"}}},
		{"InvalidBoth", "ErrorMetadata", testutils.Problems{{Message: "repeated"}, {Suggestion: "string"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
					{{.Prefix}} unreachable = 3;
				}
				message ErrorMetadata {}
				message Book {}
			`, test)
			field := f.Messages().Get(0).Fields().Get(2)
			if diff := test.problems.SetDescriptor(field).Diff(unreachableFieldType.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
