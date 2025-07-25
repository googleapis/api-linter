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

package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestMaskField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		FieldType   string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "UpdateBookRequest", "google.protobuf.FieldMask", "update_mask", nil},
		{"InvalidType", "UpdateBookRequest", "string", "update_mask", testutils.Problems{{Suggestion: "google.protobuf.FieldMask"}}},
		{"InvalidRepeated", "UpdateBookRequest", "repeated google.protobuf.FieldMask", "update_mask", testutils.Problems{{Suggestion: "google.protobuf.FieldMask"}}},
		{"IrrelevantMessage", "ModifyBookRequest", "string", "update_mask", nil},
		{"IrrelevantField", "UpdateBookRequest", "string", "modify_mask", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
				message Book {}
			`, test)
			field := file.Messages().Get(0).Fields().Get(0)
			problems := requestMaskField.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
