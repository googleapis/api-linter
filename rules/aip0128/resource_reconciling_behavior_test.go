// Copyright 2021 Google LLC
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

package aip0128

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestResourceReconcilingBehavior(t *testing.T) {
	for _, test := range []struct {
		name          string
		Style         string
		FieldName     string
		FieldBehavior string
		problems      testutils.Problems
	}{
		{"Valid", "style: DECLARATIVE_FRIENDLY", "reconciling", "[(google.api.field_behavior) = OUTPUT_ONLY]", nil},
		{"Invalid", "style: DECLARATIVE_FRIENDLY", "reconciling", "", testutils.Problems{{Message: "OUTPUT_ONLY"}}},
		{"IrrelevantStyle", "", "reconciling", "", nil},
		{"IrrelevantField", "style: DECLARATIVE_FRIENDLY", "enabled", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						{{.Style}}
					};
					bool {{.FieldName}} = 1 {{.FieldBehavior}};
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(resourceReconcilingBehavior.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
