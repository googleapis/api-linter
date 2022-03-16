// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0191

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		testName string
		filename string
		problems testutils.Problems
	}{
		{"Valid", "library.proto", testutils.Problems{}},
		{"ValidDirectory", "google/library.proto", testutils.Problems{}},
		{"ValidFileNameWithSnakeCase", "library_test.proto", testutils.Problems{}},
		{"InvalidFileNameNotSnakeCase", "library.test.proto", testutils.Problems{{Message: "invalid characters"}}},
		{"InvalidCharacterDollar", "library_$test.proto", testutils.Problems{{Message: "invalid characters"}}},
		{"InvalidCharacterSpace", "library test.proto", testutils.Problems{{Message: "invalid characters"}}},
		{"InvalidCharacterHash", "library_#test.proto", testutils.Problems{{Message: "invalid characters"}}},
		{"InvalidStable", "v1.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidBigStable", "v20.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidStableDirectory", "google/library/v1.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidAlphaUnnumbered", "v1alpha.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidAlphaNumbered", "v1alpha1.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidBetaUnnumbered", "v1beta.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidBetaNumbered", "v1beta1.proto", testutils.Problems{{Message: "proto version"}}},
		{"InvalidPoint", "v1p1beta1.proto", testutils.Problems{{Message: "proto version"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, err := builder.NewFile(test.filename).Build()
			if err != nil {
				t.Fatalf("Failed to build file.")
			}
			if diff := test.problems.SetDescriptor(f).Diff(filename.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
