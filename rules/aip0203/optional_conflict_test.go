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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestOptionalBehaviorConflict(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "Valid",
			field:    "string title = 1 [(google.api.field_behavior) = OPTIONAL];",
			problems: nil,
		},
		{
			name: "Valid",
			field: `
string title = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];`,
			problems: nil,
		},
		{
			name: "Invalid-optional-conflict",
			field: `
string title = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY,
	(google.api.field_behavior) = OPTIONAL];`,
			problems: testutils.Problems{{
				Message: "Field behavior `(google.api.field_behavior) = OPTIONAL` shouldn't be used together with other field behaviors",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
import "google/api/field_behavior.proto";
message Book {
	// Title of the book
	{{.Field}}
}`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Field string }{test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := optionalBehaviorConflict.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
