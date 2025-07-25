// Copyright 2024 Google LLC
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

package aip0217

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestReturnPartialSuccessType(t *testing.T) {
	for _, test := range []struct {
		name     string
		Type     string
		problems testutils.Problems
	}{
		{"Valid", "bool", testutils.Problems{}},
		{"InvalidType", "string", testutils.Problems{{Suggestion: "bool"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
					{{.Type}} return_partial_success = 4;
				}
			`, test)
			field := f.Messages().Get(0).Fields().ByName("return_partial_success")
			if diff := test.problems.SetDescriptor(field).Diff(returnPartialSuccessType.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
