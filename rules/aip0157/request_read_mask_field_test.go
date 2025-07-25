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

package aip0157

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestFieldMaskField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		FieldType   string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "FooRequest", "google.protobuf.FieldMask", "read_mask", nil},
		{"InvalidType", "FooRequest", "string", "read_mask", testutils.Problems{{Suggestion: "google.protobuf.FieldMask"}}},
		{"InvalidRepeated", "FooRequest", "repeated google.protobuf.FieldMask", "read_mask", testutils.Problems{{Suggestion: "google.protobuf.FieldMask"}}},
		{"IrrelevantMessage", "Foo", "string", "read_mask", nil},
		{"IrrelevantField", "FooRequest", "string", "bar", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)
			field := file.Messages().Get(0).Fields().Get(0)
			problems := requestReadMaskField.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
