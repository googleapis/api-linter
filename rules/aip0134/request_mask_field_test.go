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
	"github.com/jhump/protoreflect/desc"
)

func TestRequestMaskField(t *testing.T) {
	tests := []struct {
		testName   string
		Field      string
		descriptor func(*desc.MessageDescriptor) desc.Descriptor
		problems   testutils.Problems
	}{
		{"Valid", "google.protobuf.FieldMask update_mask = 2;", nil, testutils.Problems{}},
		{"InvalidMissing", "", nil, testutils.Problems{{Message: "`update_mask` field"}}},
		{
			"InvalidType",
			"string update_mask = 2;",
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.GetFields()[1]
			},
			testutils.Problems{{Message: "google.protobuf.FieldMask"}},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message UpdateBookRequest {
					Book book = 1;
					{{.Field}}
				}
				message Book {}
			`, test)
			message := file.GetMessageTypes()[0]
			var d desc.Descriptor = message
			if test.descriptor != nil {
				d = test.descriptor(message)
			}
			if diff := test.problems.SetDescriptor(d).Diff(requestMaskField.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
