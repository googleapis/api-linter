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

func TestRequestParentReference(t *testing.T) {
	t.Run("Present", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message CreateBookRequest {
				string parent = 1 [(google.api.resource_reference) = {
					type: "library.googleapis.com/Publisher"
				}];
			}
		`)
		if diff := (testutils.Problems{}).Diff(requestParentReference.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	})
	t.Run("Absent", func(t *testing.T) {
		for _, test := range []struct {
			name      string
			FieldName string
			problems  testutils.Problems
		}{
			{"Error", "parent", testutils.Problems{{Message: "google.api.resource_reference"}}},
			{"Irrelevant", "something_else", testutils.Problems{}},
		} {
			t.Run(test.name, func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					import "google/api/resource.proto";
					message CreateBookRequest {
						string {{.FieldName}} = 1;
					}
				`, test)
				field := f.GetMessageTypes()[0].GetFields()[0]
				if diff := test.problems.SetDescriptor(field).Diff(requestParentReference.Lint(f)); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})
}
