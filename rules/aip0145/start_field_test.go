// Copyright 2022 Google LLC
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

package aip0145

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestStartField(t *testing.T) {
	for _, test := range []struct {
		name     string
		Suffix   string
		problems testutils.Problems
	}{
		{"Valid", "_page", nil},
		{"Invalid", "", testutils.Problems{{Message: "start_xxx"}}},
		{"InvalidTrailingUnderscore", "_", testutils.Problems{{Message: "start_xxx"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Foo {
					int32 start{{ .Suffix }} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(startField.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
