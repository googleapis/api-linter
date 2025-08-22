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

package aip4232

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestRequiredFields(t *testing.T) {
	tests := []struct {
		testName      string
		Signature     string
		FieldBehavior string
		problems      testutils.Problems
	}{
		{"Valid", "name,paperback_only,shelf", "REQUIRED", nil},
		{"ValidNested", "name,paperback_only,shelf.name", "REQUIRED", nil},
		{"ValidNotRequired", "paperback_only", "OPTIONAL", nil},
		{"Invalid", "paperback_only", "REQUIRED", testutils.Problems{{Message: "missing at least one"}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/client.proto";
				import "google/api/field_behavior.proto";
				service Library {
					rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResponse) {
						option (google.api.method_signature) = "{{.Signature}}";
					}
				}
				message ArchiveBookRequest {
					string name = 1 [(google.api.field_behavior) = {{.FieldBehavior}}];

					bool paperback_only = 2;

					Shelf shelf = 3 [(google.api.field_behavior) = {{.FieldBehavior}}];
				}
				message ArchiveBookResponse {}
				message Shelf {
					string name = 1;
				}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(method).Diff(requiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
