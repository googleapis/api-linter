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
		{"InvalidFileNameNotSnakeCase", "library.test.proto", testutils.Problems{{Message: "The filename has invalid characters."}}},
		{"InvalidCharacterUpperCase", "library_Test.proto", testutils.Problems{{Message: "The filename has invalid characters."}}},
		{"InvalidCharacterDollar", "library_$test.proto", testutils.Problems{{Message: "The filename has invalid characters."}}},
		{"InvalidCharacterSpace", "library test.proto", testutils.Problems{{Message: "The filename has invalid characters."}}},
		{"InvalidCharacterHash", "library_#test.proto", testutils.Problems{{Message: "The filename has invalid characters."}}},
		{"InvalidStable", "v1.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidBigStable", "v20.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidStableDirectory", "google/library/v1.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidAlphaUnnumbered", "v1alpha.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidAlphaNumbered", "v1alpha1.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidBetaUnnumbered", "v1beta.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidBetaNumbered", "v1beta1.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
		{"InvalidPoint", "v1p1beta1.proto", testutils.Problems{{Message: "The proto version must not be used as the filename."}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			files := testutils.ParseProtoStrings(t, map[string]string{
				test.filename: "",
			})
			f := files[test.filename]
			if diff := test.problems.SetDescriptor(f).Diff(filename.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
