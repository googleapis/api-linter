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

func TestOptionalBehaviorConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		problems testutils.Problems
	}{
		{
			name: "Valid-NoneOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3;

string author = 4;`,
			problems: nil,
		},
		{
			name: "Valid-AllOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3 [(google.api.field_behavior) = OPTIONAL];

string author = 4 [(google.api.field_behavior) = OPTIONAL];`,
			problems: nil,
		},
		{
			name: "Invalid-PartialOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3 [(google.api.field_behavior) = OPTIONAL];

string author = 4;`,
			problems: testutils.Problems{{
				Message: "Within a single message, either all optional fields should be indicated, or none of them should be.",
			}},
		},
		{
			name: "Valid-IgnoreOneofFields",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3 [(google.api.field_behavior) = OPTIONAL];

oneof other {
	string author = 4;	
}`,
			problems: nil,
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
			// author field in the test will get the warning
			f := file.GetMessageTypes()[0].GetFields()[3]
			problems := optionalBehaviorConsistency.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
