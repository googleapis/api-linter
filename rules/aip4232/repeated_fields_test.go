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

func TestRepeatedFields(t *testing.T) {
	tests := []struct {
		testName  string
		FirstSig  string
		SecondSig string
		problems  testutils.Problems
	}{
		{"Valid", "name,paperback_only,book.editions", "name,editions", nil},
		{"InvalidFirstSignature", "name,book.shelves.shelf,paperback_only", "name,editions", testutils.Problems{{Message: "only the last"}}},
		{"InvalidSecondSignature", "name,book.editions", "name,book.shelves.shelf,paperback_only", testutils.Problems{{Message: "only the last"}}},
		{"BothInvalid", "name,book.shelves.shelf,paperback_only", "name,book.shelves.shelf,paperback_only", testutils.Problems{{Message: "only the last"}, {Message: "only the last"}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/client.proto";
				service Library {
					rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResponse) {
						option (google.api.method_signature) = "{{.FirstSig}}";
						option (google.api.method_signature) = "{{.SecondSig}}";
					}
				}
				message ArchiveBookRequest {
					string name = 1;

					bool paperback_only = 2;

					Book book = 3;
				}
				message Book {
					string name = 1;

					repeated int32 editions = 2;

					repeated Shelf shelves = 3;
				}
				message Shelf {
					string shelf = 1;
				}
				message ArchiveBookResponse {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(method).Diff(repeatedFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
