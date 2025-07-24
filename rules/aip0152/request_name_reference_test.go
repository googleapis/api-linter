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

package aip0152

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestNameReference(t *testing.T) {
	t.Run("Present", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message RunWriteBookJobRequest {
				string name = 1 [(google.api.resource_reference) = {
					type: "library.googleapis.com/Book"
				}];
			}
		`)
		if diff := (testutils.Problems{}).Diff(requestNameReference.Lint(f)); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("Absent", func(t *testing.T) {
		for _, test := range []struct {
			name      string
			FieldName string
			problems  testutils.Problems
		}{
			{"Error", "name", testutils.Problems{{Message: "google.api.resource_reference"}}},
			{"Irrelevant", "something_else", nil},
		} {
			t.Run(test.name, func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					import "google/api/resource.proto";
					message RunWriteBookJobRequest {
						string {{.FieldName}} = 1;
					}
				`, test)
				field := f.Messages()[0].Fields()[0]
				if diff := test.problems.SetDescriptor(field).Diff(requestNameReference.Lint(f)); diff != "" {
					t.Error(diff)
				}
			})
		}
	})
}
