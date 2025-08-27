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

package aip0162

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestDeleteRevisionResponseMessageName(t *testing.T) {
	for _, test := range []struct {
		name         string
		MethodName   string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "DeleteBookRevision", "Book", nil},
		{"Invalid", "DeleteBookRevision", "DeleteBookRevisionResponse", testutils.Problems{{Suggestion: "Book"}}},
		{"Irrelevant", "DeleteBook", "DeleteBookResponse", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.ResponseType}});
				}
				message {{.MethodName}}Request {}
				message {{.ResponseType}} {}
			`, test)

			method := file.Services().Get(0).Methods().Get(0)
			problems := deleteRevisionResponseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeleteRevisionResponseMessageNameLRO(t *testing.T) {
	for _, test := range []struct {
		name         string
		MethodName   string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "DeleteBookRevision", "Book", nil},
		{"Invalid", "DeleteBookRevision", "DeleteBookRevisionResponse", testutils.Problems{{Message: "Book"}}},
		{"Irrelevant", "DeleteBook", "DeleteBookResponse", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (google.longrunning.Operation) {
					  option (google.longrunning.operation_info) = {
				        response_type: "{{.ResponseType}}"
						metadata_type: "OperationMetadata"
					  };
					};
				}
				message {{.MethodName}}Request {}
				message {{.ResponseType}} {}
				message OperationMetadata {}
			`, test)

			method := file.Services().Get(0).Methods().Get(0)
			problems := deleteRevisionResponseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
