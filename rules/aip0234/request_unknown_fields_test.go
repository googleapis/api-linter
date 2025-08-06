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

package aip0234

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestRequestUnknownFields(t *testing.T) {
	for _, test := range []struct {
		name        string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid-AllowMissing", "BatchUpdateBooks", "allow_missing", testutils.Problems{}},
		{"Valid-Parent", "BatchUpdateBooks", "parent", testutils.Problems{}},
		{"Valid-RequestID", "BatchCreateBooks", "request_id", testutils.Problems{}},
		{"Valid-Requests", "BatchUpdateBooks", "requests", testutils.Problems{}},
		{"Valid-UpdateMask", "BatchUpdateBooks", "update_mask", testutils.Problems{}},
		{"Valid-ValidateOnly", "BatchUpdateBooks", "validate_only", testutils.Problems{}},
		{"Invalid", "BatchUpdateBooks", "foo", testutils.Problems{{Message: "Unexpected field"}}},
		{"IrrelevantMessage", "UpdateBooks", "foo", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}}Request {
					repeated string {{.FieldName}} = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(requestUnknownFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
