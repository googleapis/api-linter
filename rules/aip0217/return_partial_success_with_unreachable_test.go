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

func TestReturnPartialSuccessWithUnreachable(t *testing.T) {
	for _, test := range []struct {
		name          string
		ResponseField string
		problems      testutils.Problems
	}{
		{"Valid", "repeated string unreachable = 1;", testutils.Problems{}},
		{"InvalidMissingUnreachable", "", testutils.Problems{{Message: "repeated string unreachable"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc ListBooks(ListBooksRequest) returns (ListBooksResponse);
				}

				message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
					bool return_partial_success = 4;
				}

				message ListBooksResponse {
					{{.ResponseField}}
				}
			`, test)
			field := f.GetMessageTypes()[0].FindFieldByName("return_partial_success")
			if diff := test.problems.SetDescriptor(field).Diff(returnPartialSuccessWithUnreachable.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
