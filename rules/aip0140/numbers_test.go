// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestNumbers(t *testing.T) {
	for _, test := range []struct {
		Name     string
		problems testutils.Problems
	}{
		{"foo", nil},
		{"foo_bar", nil},
		{"foo_123", testutils.Problems{{Message: "number"}}},
		{"bar_5foo", testutils.Problems{{Message: "number"}}},
		// We do not need to test field names that start with a number
		// (e.g. 123_foo) because it is invalid syntax in protocol buffers.
	} {
		t.Run(test.Name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message Baz {
					string {{.Name}} = 1;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(numbers.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
