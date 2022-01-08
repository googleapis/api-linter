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

package aip0143

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestFieldNames(t *testing.T) {
	tests := []struct {
		FieldName string
		problems  testutils.Problems
	}{
		{"something_random", testutils.Problems{}},
		{"content_type", testutils.Problems{{Suggestion: "mime_type"}}},
		{"country", testutils.Problems{{Suggestion: "region_code"}}},
		{"country_code", testutils.Problems{{Suggestion: "region_code"}}},
		{"region_code", testutils.Problems{}},
		{"currency", testutils.Problems{{Suggestion: "currency_code"}}},
		{"currency_code", testutils.Problems{}},
		{"language", testutils.Problems{{Suggestion: "language_code"}}},
		{"language_code", testutils.Problems{}},
		{"mime", testutils.Problems{{Suggestion: "mime_type"}}},
		{"mimetype", testutils.Problems{{Suggestion: "mime_type"}}},
		{"mime_type", testutils.Problems{}},
		{"timezone", testutils.Problems{{Suggestion: "time_zone"}}},
		{"time_zone", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.FieldName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message Foo {
					string {{.FieldName}} = 1;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			problems := fieldNames.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
