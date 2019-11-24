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

func TestDefinedInFile(t *testing.T) {
	tests := []struct {
		testName     string
		Package      string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "", "WriteBookResponse", testutils.Problems{}},
		{"ValidPkg", "package test;", "WriteBookResponse", testutils.Problems{}},
		{"InvalidTypo", "", "WriteBookRepsonse", testutils.Problems{{Message: "should be defined"}}},
		{"InvalidTypoPkg", "package test;", "WriteBookRepsonse", testutils.Problems{{Message: "should be defined"}}},
		{"ValidExternal", "", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidExternalPkg", "package test;", "google.protobuf.Empty", testutils.Problems{}},
	}
	for _, test := range tests {
		f := testutils.ParseProto3Tmpl(t, `
			{{.Package}}
			import "google/longrunning/operations.proto";
			service Library {
				rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "{{.ResponseType}}"
						metadata_type: "OperationMetadata"
					};
				}
			}
			message WriteBookRequest {}
			message WriteBookResponse {}
			message OperationMetadata {}
		`, test)
		m := f.GetServices()[0].GetMethods()[0]
		problems := lroDefinedInFile.Lint(m)
		if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
			t.Errorf(diff)
		}
	}
}
