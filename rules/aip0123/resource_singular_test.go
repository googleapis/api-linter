// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
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

func TestResourceSingular(t *testing.T) {
	for _, test := range []struct {
		name     string
		Singular string
		problems testutils.Problems
	}{
		{
			"Valid",
			`singular: "book"`,
			nil,
		},
		{
			"InvalidDoesntMatchType",
			`singular: "shelf"`,
			testutils.Problems{{
				Message: "book",
			}},
		},
		{
			"InvalidMissing",
			``,
			testutils.Problems{{
				Message: "Resources should declare singular",
			}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
					{{.Singular}}
				};
				string name = 1;
			}
			`, test)
			m := f.Messages().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(resourceSingular.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
