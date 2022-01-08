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

package aip0131

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestRequestNameBehavior(t *testing.T) {
	for _, test := range []struct {
		name          string
		FieldName     string
		FieldBehavior string
		problems      testutils.Problems
	}{
		{"Valid", "name", " [(google.api.field_behavior) = REQUIRED]", testutils.Problems{}},
		{"Missing", "name", "", testutils.Problems{{Message: "(google.api.field_behavior) = REQUIRED"}}},
		{"Irrelevant", "something_else", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				message GetBookRequest {
					string {{.FieldName}} = 1{{.FieldBehavior}};
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(requestNameBehavior.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
