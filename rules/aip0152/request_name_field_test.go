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

package aip0152

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRequestNameField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		Field       string
		problems    testutils.Problems
	}{
		{"Valid", "RunWriteBookJobRequest", `string name = 1;`, nil},
		{"InvalidMissing", "RunWriteBookJobRequest", "", testutils.Problems{{Message: "has no"}}},
		{"InvalidType", "RunWriteBookJobRequest", "bytes name = 1;", testutils.Problems{{Suggestion: "string"}}},
		{"InvalidRepeated", "RunWriteBookJobRequest", "repeated string name = 1;", testutils.Problems{{Suggestion: "string"}}},
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
			var d protoreflect.Descriptor = f.Messages().Get(0)
			if test.name == "InvalidType" || test.name == "InvalidRepeated" {
				d = f.Messages().Get(0).Fields().Get(0)
			}
			problems := requestNameField.Lint(f)
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}
