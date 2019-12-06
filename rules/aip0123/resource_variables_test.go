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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceVariables(t *testing.T) {
	for _, test := range []struct {
		name     string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "publishers/{publisher}/electronicBooks/{electronic_book}", testutils.Problems{}},
		{"CamelCase", "publishers/{publisher}/electronicBooks/{electronicBook}", testutils.Problems{{
			Message: "publishers/{publisher}/electronicBooks/{electronic_book}",
		}}},
		{"ID", "publishers/{publisher}/electronicBooks/{electronic_book_id}", testutils.Problems{{
			Message: "publishers/{publisher}/electronicBooks/{electronic_book}",
		}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "{{.Pattern}}"
					};
					string name = 1;
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(resourceVariables.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
