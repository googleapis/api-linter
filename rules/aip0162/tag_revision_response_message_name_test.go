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

func TestTagRevisionResponseMessageName(t *testing.T) {
	for _, test := range []struct {
		name         string
		Method       string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "TagBookRevision", "Book", nil},
		{"Invalid", "TagBookRevision", "TagBookRevisionResponse", testutils.Problems{{Suggestion: "Book"}}},
		{"Irrelevant", "AcquireBook", "PurgeBooksResponse", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.Method}}(TagBookRevisionRequest) returns ({{.ResponseType}});
				}
				message TagBookRevisionRequest {}
				message {{.ResponseType}} {}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(tagRevisionResponseMessageName.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestTagRevisionResponseMessageNameLRO(t *testing.T) {
	for _, test := range []struct {
		name         string
		Method       string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "TagBookRevision", "Book", nil},
		{"Invalid", "TagBookRevision", "TagBookRevisionResponse", testutils.Problems{{Message: "Book"}}},
		{"Irrelevant", "AcquireBook", "PurgeBooksResponse", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.Method}}(TagBookRevisionRequest) returns (google.longrunning.Operation) {
					  option (google.longrunning.operation_info) = {
					    response_type: "{{.ResponseType}}"
						metadata_type: "OperationMetadata"
					  };
					};
				}
				message TagBookRevisionRequest {}
				message {{.ResponseType}} {}
				message OperationMetadata {}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(tagRevisionResponseMessageName.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
