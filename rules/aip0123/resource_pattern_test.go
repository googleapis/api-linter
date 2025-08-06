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

package aip0123

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestResourcePattern(t *testing.T) {
	for _, test := range []struct {
		name     string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", `pattern: "publishers/{publisher}/books/{book}"`, testutils.Problems{}},
		{"ValidCamel", `pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"`, testutils.Problems{}},
		{"Missing", "", testutils.Problems{{Message: "declare resource name pattern"}}},
		{"SnakeCase", `pattern: "book_publishers/{book_publisher}/books/{book}"`, testutils.Problems{{
			Message: "bookPublishers/{book_publisher}/books/{book}",
		}}},
		{"HasSpaces", `pattern: "publishers/{publisher}/ books /{book}"`, testutils.Problems{{
			Message: "Resource patterns should not have spaces",
		}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						{{.Pattern}}
					};
					string name = 1;
				}
			`, test)
			m := f.Messages().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(resourcePattern.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
