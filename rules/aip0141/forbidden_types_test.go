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

package aip0141

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestForbiddenTypes(t *testing.T) {
	tests := []struct {
		TypeName string
		problems testutils.Problems
	}{
		{"int32", testutils.Problems{}},
		{"int64", testutils.Problems{}},
		{"uint32", testutils.Problems{{Suggestion: "int32"}}},
		{"uint64", testutils.Problems{{Suggestion: "int64"}}},
		{"fixed32", testutils.Problems{{Suggestion: "int32"}}},
		{"fixed64", testutils.Problems{{Suggestion: "int64"}}},
	}
	for _, test := range tests {
		t.Run(test.TypeName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message Book {
					{{.TypeName}} pages = 1;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			problems := forbiddenTypes.Lint(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
