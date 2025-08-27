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

package aip0123

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestResourceDefinitionPattern(t *testing.T) {
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
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				option (google.api.resource_definition) = {
					type: "library.googleapis.com/Book"
					{{.Pattern}}
				};
			`, test)
			got := resourceDefinitionPatterns.Lint(f)
			if diff := test.problems.SetDescriptor(f).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
