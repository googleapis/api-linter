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

package aip0151

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestLROMetadata(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		Metadata   string
		problems   testutils.Problems
	}{
		{"Valid", "WriteBook", "WriteBookMetadata", testutils.Problems{}},
		{"InvalidEmptyString", "WriteBook", "", testutils.Problems{{Message: "must set the metadata type"}}},
		{"InvalidGPEmpty", "WriteBook", "google.protobuf.Empty", testutils.Problems{{Message: "Empty"}}},
		{"InvalidGPEmptyDelete", "DeleteBook", "google.protobuf.Empty", testutils.Problems{{Message: "Empty"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.MethodName}}Response"
							metadata_type: "{{.Metadata}}"
						};
					}
				}
				message {{.MethodName}}Request {}
			`, test)
			method := f.GetServices()[0].GetMethods()[0]
			problems := lroMetadata.Lint(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
