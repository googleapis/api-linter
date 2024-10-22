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

package aip0128

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestResourceReconcilingField(t *testing.T) {
	for _, test := range []struct {
		name        string
		Style       string
		Reconciling string
		problems    testutils.Problems
	}{
		{"ValidMissingNotDF", "", "", nil},
		{"ValidPresentNotDF", "", "bool reconciling = 2;", nil},
		{"ValidBadTypeNotDF", "", "int32 reconciling = 2;", nil},
		{"ValidPresentDF", "style: DECLARATIVE_FRIENDLY", "bool reconciling = 2;", nil},
		{"InvalidMissingDF", "style: DECLARATIVE_FRIENDLY", "", testutils.Problems{{Message: "reconciling"}}},
		{"InvalidBadTypeDF", "style: DECLARATIVE_FRIENDLY", "int32 reconciling = 2;", testutils.Problems{{Suggestion: "bool"}}},
		{"InvalidRepeatedDF", "style: DECLARATIVE_FRIENDLY", "repeated bool reconciling = 2;", testutils.Problems{{Suggestion: "bool"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						{{.Style}}
					};
					string name = 1;
					{{.Reconciling}}
				}
				message DeleteBookRequest {}
			`, test)
			var d desc.Descriptor = f.GetMessageTypes()[0]
			if test.name == "InvalidBadTypeDF" || test.name == "InvalidRepeatedDF" {
				d = f.GetMessageTypes()[0].GetFields()[1]
			}
			if diff := test.problems.SetDescriptor(d).Diff(resourceReconcilingField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
