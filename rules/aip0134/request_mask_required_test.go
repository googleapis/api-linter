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

func TestRequestMaskFieldRequired(t *testing.T) {
	tests := []struct {
		name     string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "google.protobuf.FieldMask update_mask = 2;", nil},
		{"InvalidMissing", "", testutils.Problems{{Message: "`update_mask` field"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message UpdateBookRequest {
					Book book = 1;
					{{.Field}}
				}
				message Book {}
			`, test)
			message := file.GetMessageTypes()[0]
			problems := requestMaskRequired.Lint(file)
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
