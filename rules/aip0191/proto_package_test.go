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

package aip0191

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestProtoPkg(t *testing.T) {
	tests := []struct {
		testName string
		filename string
		Pkg      string
		problems testutils.Problems
	}{
		{"Valid", "google/example/library/v1/library.proto", "google.example.library.v1", testutils.Problems{}},
		{"InvalidPackage", "google/example/library/v1/library.proto", "google.library.v1", testutils.Problems{{Message: "directory structure"}}},
		{"InvalidDirectory", "google/v1/library.proto", "google.example.library.v1", testutils.Problems{{Message: "directory structure"}}},
		{"IgnoreRootDirectory", "library.proto", "google.example.library.v1", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			files := testutils.ParseProto3Tmpls(t, map[string]string{
				test.filename: `package {{.Pkg}};`,
			}, test)
			f := files[test.filename]
			if diff := test.problems.SetDescriptor(f).Diff(protoPkg.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
