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

package aip0122

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNameSuffix(t *testing.T) {
	for _, test := range []struct {
		name      string
		FieldName string
		problems  testutils.Problems
	}{
		{"Valid", "publisher", testutils.Problems{}},
		{"ValidStandardDisplay", "display_name", testutils.Problems{}},
		{"ValidStandardGiven", "given_name", testutils.Problems{}},
		{"ValidStandardFamily", "family_name", testutils.Problems{}},
		{"ValidStandardFull", "full_resource_name", testutils.Problems{}},
		{"SkipValidDisplayNameSuffix", "foo_display_name", testutils.Problems{}},
		{"Invalid", "author_name", testutils.Problems{{Suggestion: "author"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			message Book {
				string name = 1;
				string {{.FieldName}} = 2;
			}
		`, test)
		field := f.GetMessageTypes()[0].GetFields()[1]
		if diff := test.problems.SetDescriptor(field).Diff(nameSuffix.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	}
}
