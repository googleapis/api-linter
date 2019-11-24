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

package aip0192

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

// TODO: add more tests for different descriptors.
func TestHasComments(t *testing.T) {
	file := testutils.ParseProto3String(t, `
		// This is a book.
		message Book {
			// The resource name.
			string name = 1;
			string title = 2;
		}
	`)
	field := file.GetMessageTypes()[0].GetFields()[1]
	wantProblems := testutils.Problems{{Descriptor: field}}
	if diff := wantProblems.Diff(hasComments.Lint(field)); diff != "" {
		t.Errorf(diff)
	}
}
