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

func TestNameNeverOptional(t *testing.T) {
	for _, test := range []struct {
		name      string
		FieldName string
		NameField string
		Label     string
		problems  testutils.Problems
	}{
		{"Valid", "name", "", "", testutils.Problems{}},
		{"ValidAlternativeName", "resource", "resource", "", testutils.Problems{}},
		{"InvalidProto3Optional", "name", "", "optional", testutils.Problems{{Message: "never be labeled"}}},
		{"SkipNameFieldDNE", "name", "does_not_exist", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						name_field: "{{.NameField}}"
					};

					{{.Label}} string {{.FieldName}} = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(nameNeverOptional.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
