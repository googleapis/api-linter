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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestUnorderedList(t *testing.T) {
	testCases := []struct {
		name     string
		Field    string
		problems testutils.Problems
	}{
		{
			name:     "ValidSingularNoAnnotation",
			Field:    `string title = 1;`,
			problems: nil,
		},
		{
			name:     "ValidRepeatedNoAnnotation",
			Field:    `repeated string authors = 1;`,
			problems: nil,
		},
		{
			name:     "ValidRepeatedAnnotation",
			Field:    `repeated string authors = 1 [(google.api.field_behavior) = UNORDERED_LIST];`,
			problems: nil,
		},
		{
			name:     "InvalidSingularAnnotation",
			Field:    `string title = 1 [(google.api.field_behavior) = UNORDERED_LIST];`,
			problems: testutils.Problems{{Message: "must not be applied to non-repeated fields"}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				message Book {
					{{.Field}}
				}`, test)
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := unorderedListRepeated.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
