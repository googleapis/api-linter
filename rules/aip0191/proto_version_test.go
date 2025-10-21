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

package aip0191

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestSyntax(t *testing.T) {
	// Set up the two permutations.
	tests := []struct {
		testName string
		isProto3 bool
		problems testutils.Problems
	}{
		{"Valid", true, testutils.Problems{}},
		{"Invalid", false, testutils.Problems{{Suggestion: `syntax = "proto3";`}}},
	}

	// Run each permutation as an individual test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Build an appropriate file descriptor.
			f, err := builder.NewFile("library.proto").SetProto3(test.isProto3).Build()
			if err != nil {
				t.Fatalf("Could not build file descriptor.")
			}

			// Lint the file, and ensure we got the expected problems.
			if diff := test.problems.SetDescriptor(f).Diff(syntax.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
