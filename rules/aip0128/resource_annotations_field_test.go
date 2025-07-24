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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResourceAnnotationsField(t *testing.T) {
	for _, test := range []struct {
		name        string
		Style       string
		Annotations string
		problems    testutils.Problems
	}{
		{"ValidMissingNotDF", "", "", nil},
		{"ValidPresentNotDF", "", "map<string, string> annotations = 2;", nil},
		{"ValidBadTypeNotDF", "", "int32 annotations = 2;", nil},
		{"ValidPresentDF", "style: DECLARATIVE_FRIENDLY", "map<string, string> annotations = 2;", nil},
		{"InvalidMissingDF", "style: DECLARATIVE_FRIENDLY", "", testutils.Problems{{Message: "annotations"}}},
		{"InvalidBadTypeDF", "style: DECLARATIVE_FRIENDLY", "map<string, int32> annotations = 2;", testutils.Problems{{Suggestion: "map<string, string>"}}},
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
					{{.Annotations}}
				}
				message DeleteBookRequest {}
			`, test)
			var d protoreflect.Descriptor = f.Messages()[0]
			if test.name == "InvalidBadTypeDF" {
				d = f.Messages()[0].Fields()[1]
			}
			if diff := test.problems.SetDescriptor(d).Diff(resourceAnnotationsField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
