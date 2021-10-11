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

package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestNameReferenceType(t *testing.T) {
	for _, test := range []struct {
		testName  string
		Reference string
		problems  testutils.Problems
	}{
		{"Valid", "type", nil},
		{"Invalid", "child_type", testutils.Problems{{Message: "should be a direct"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message GetBookRequest {
					string name = 1 [(google.api.resource_reference).{{.Reference}} = "library.googleapis.com/Book"];
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(requestNameReferenceType.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
