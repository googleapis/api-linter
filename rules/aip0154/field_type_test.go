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

package aip0154

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestFieldType(t *testing.T) {
	for _, test := range []struct {
		name     string
		Type     string
		problems testutils.Problems
	}{
		{"ValidString", "string", nil},
		{"Invalid", "bytes", testutils.Problems{{Suggestion: "string"}}},
		{"InvalidRepeated", "repeated string", testutils.Problems{{Suggestion: "string"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Book {
					{{.Type}} etag = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(fieldType.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
