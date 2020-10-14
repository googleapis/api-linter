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

package aip0164

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRequestNameField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		Field       string
		problems    testutils.Problems
	}{
		{"Valid", "UndeleteBookRequest", `string name = 1;`, nil},
		{"InvalidMissing", "UndeleteBookRequest", "", testutils.Problems{{Message: "has no"}}},
		{"InvalidType", "UndeleteBookRequest", "bytes name = 1;", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantMessage", "RemoveBookRequest", "", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.Field}}
					bytes other_field = 2;
				}
			`, test)
			var d desc.Descriptor = f.GetMessageTypes()[0]
			if test.name == "InvalidType" {
				d = f.GetMessageTypes()[0].GetFields()[0]
			}
			problems := requestNameField.Lint(f)
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}
