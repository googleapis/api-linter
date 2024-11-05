// Copyright 2024 Google LLC
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

func TestUpdateMaskOptionalBehavior(t *testing.T) {
	tests := []struct {
		name     string
		Behavior string
		problems testutils.Problems
	}{
		{"Valid", "[(google.api.field_behavior) = OPTIONAL]", nil},
		{"InvalidWrong", "[(google.api.field_behavior) = REQUIRED]", testutils.Problems{{Message: "must have `OPTIONAL`"}}},
		{"InvalidMissing", "", testutils.Problems{{Message: "must have `OPTIONAL`"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/protobuf/field_mask.proto";
				message UpdateBookRequest {
					Book book = 1;
					google.protobuf.FieldMask update_mask = 2 {{.Behavior}};
				}
				message Book {}
			`, test)
			field := file.GetMessageTypes()[0].FindFieldByName("update_mask")
			problems := updateMaskOptionalBehavior.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
