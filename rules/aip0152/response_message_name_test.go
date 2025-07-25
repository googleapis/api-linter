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

package aip0152

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseMessageName(t *testing.T) {
	for _, test := range []struct {
		name         string
		Method       string
		ResponseType string
		LROType      string
		problems     testutils.Problems
	}{
		{"Valid", "RunWriteBookJob", "google.longrunning.Operation", "RunWriteBookJobResponse", nil},
		{"InvalidNotLRO", "RunWriteBookJob", "RunWriteBookJobResponse", "", testutils.Problems{{Suggestion: "google.longrunning.Operation"}}},
		{"InvalidLROType", "RunWriteBookJob", "google.longrunning.Operation", "RunWriteBookJobRequest", testutils.Problems{{Suggestion: "RunWriteBookJobResponse"}}},
		{"Irrelevant", "AcquireBook", "RunWriteBookJobResponse", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.Method}}(RunWriteBookJobRequest) returns ({{.ResponseType}}) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.LROType}}"
						};
					}
				}
				message RunWriteBookJobRequest {}
				message RunWriteBookJobResponse {}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(responseMessageName.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
