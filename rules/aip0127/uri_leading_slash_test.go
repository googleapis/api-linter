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

package aip0127

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestURILeadingSlash(t *testing.T) {
	for _, test := range []struct {
		name     string
		URI      string
		problems testutils.Problems
	}{
		{"Valid", "/v1/{name=publishers/*/books/*}", nil},
		{"Invalid", "v1/{name=publishers/*/books/*}", testutils.Problems{{Message: "leading slash"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/api/annotations.proto";
			import "google/protobuf/empty.proto";

			service Foo {
				rpc FrobFoo(google.protobuf.Empty) returns (google.protobuf.Empty) {
					option (google.api.http) = {
						post: "{{.URI}}"
						body: "*"
					};
				}
			}
		`, test)
		m := f.Services().Get(0).Methods().Get(0)
		if diff := test.problems.SetDescriptor(m).Diff(leadingSlash.Lint(f)); diff != "" {
			t.Error(diff)
		}
	}
}
