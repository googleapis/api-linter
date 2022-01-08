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

package aip0165

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestResponseMessageName(t *testing.T) {
	for _, test := range []struct {
		name         string
		Method       string
		ResponseType string
		LROType      string
		problems     testutils.Problems
	}{
		{"Valid", "PurgeBooks", "google.longrunning.Operation", "PurgeBooksResponse", nil},
		{"InvalidNotLRO", "PurgeBooks", "PurgeBooksResponse", "", testutils.Problems{{Suggestion: "google.longrunning.Operation"}}},
		{"InvalidLROType", "PurgeBooks", "google.longrunning.Operation", "PurgeBooksRequest", testutils.Problems{{Suggestion: "PurgeBooksResponse"}}},
		{"Irrelevant", "AcquireBook", "PurgeBooksResponse", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.Method}}(PurgeBooksRequest) returns ({{.ResponseType}}) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.LROType}}"
						};
					}
				}
				message PurgeBooksRequest {}
				message PurgeBooksResponse {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(responseMessageName.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
