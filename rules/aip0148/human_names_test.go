// Copyright 2020 Google LLC
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

package aip0148

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHumanNames(t *testing.T) {
	for _, test := range []struct {
		FieldName string
		problems  testutils.Problems
	}{
		{"given_name", nil},
		{"family_name", nil},
		{"first_name", testutils.Problems{{Suggestion: "given_name"}}},
		{"last_name", testutils.Problems{{Suggestion: "family_name"}}},
	} {
		t.Run(test.FieldName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Person {
					string {{.FieldName}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(humanNames.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
