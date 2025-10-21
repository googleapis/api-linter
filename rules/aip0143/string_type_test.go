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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestFieldTypes(t *testing.T) {
	tests := []struct {
		testName  string
		FieldType string
		FieldName string
		problems  testutils.Problems
	}{
		{"Irrelevant", "int32", "book_count", testutils.Problems{}},
		{"Valid", "string", "language_code", testutils.Problems{}},
		{"ValidTimeZone", "google.type.TimeZone", "time_zone", nil},
		{"InvalidScalar", "bytes", "language_code", testutils.Problems{{Suggestion: "string"}}},
		{"InvalidEnum", "Language", "language_code", testutils.Problems{{Suggestion: "string"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/type/datetime.proto";

				message Foo {
					{{.FieldType}} {{.FieldName}} = 1;
				}
				enum Language {
					LANGUAGE_UNSPECIFIED = 0;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			problems := fieldTypes.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
