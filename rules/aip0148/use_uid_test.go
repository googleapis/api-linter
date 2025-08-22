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

package aip0148

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestUseUid(t *testing.T) {
	for _, test := range []struct {
		FieldName string
		problems  testutils.Problems
	}{
		{"uid", nil},
		{"id", testutils.Problems{{Suggestion: "uid"}}},
		{"foo_id", nil},
	} {
		t.Run(test.FieldName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				message Person {
					option (google.api.resource) = {
						type: "foo.googleapis.com/Person"
						pattern: "people/{person}"
					};
					string name = 1;

					string {{.FieldName}} = 2;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(1)
			if diff := test.problems.SetDescriptor(field).Diff(useUid.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
