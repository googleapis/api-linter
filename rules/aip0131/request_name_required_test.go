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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestHasNameField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "GetBookRequest", "name", testutils.Problems{}},
		{"InvalidName", "GetBookRequest", "id", testutils.Problems{{Message: "name"}}},
		{"Irrelevant", "AcquireBookRequest", "id", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			message {{.MessageName}} {
				string {{.FieldName}} = 1;
			}`, test)

			// Run the lint rule, and establish that it returns the correct problems.
			problems := requestNameRequired.Lint(f)
			if diff := test.problems.SetDescriptor(f.Messages()[0]).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}
