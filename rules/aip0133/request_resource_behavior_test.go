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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestResourceBehavior(t *testing.T) {
	for _, test := range []struct {
		name          string
		MessageName   string
		FieldType     string
		FieldBehavior string
		problems      testutils.Problems
	}{
		{"Valid", "CreateBookRequest", "Book", " [(google.api.field_behavior) = REQUIRED]", testutils.Problems{}},
		{"Missing", "CreateBookRequest", "Book", "", testutils.Problems{{Message: "(google.api.field_behavior) = REQUIRED"}}},
		{"IrrelevantOtherField", "CreateBookRequest", "Parchment", "", testutils.Problems{}},
		{"IrrelevantNotCreateRequest", "FrobBookRequest", "Book", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				message {{.MessageName}} {
					{{.FieldType}} book = 1{{.FieldBehavior}};
				}
				message {{.FieldType}} {}
			`, test)
			field := f.Messages()[0].Fields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(requestResourceBehavior.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
