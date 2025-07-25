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

package aip0136

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestDeclarativeFriendly(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		ImpOnly    string
		problems   testutils.Problems
	}{
		{"ValidGet", "GetBook", "", nil},
		{"ValidList", "ListBooks", "", nil},
		{"ValidCreate", "CreateBook", "", nil},
		{"ValidUpdate", "UpdateBook", "", nil},
		{"ValidDelete", "DeleteBook", "", nil},
		{"ValidUndelete", "UndeleteBook", "", nil},
		{"ValidBatch", "BatchGetBooks", "", nil},
		{"ValidCustomImperativeOnly", "FrobBook", "IMPERATIVE ONLY.", nil},
		{"Invalid", "FrobBook", "", testutils.Problems{{Message: "avoid custom methods"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				service Library {
					// The {{.MethodName}} method.
					// (-- {{.ImpOnly}} --)
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response);
				}

				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						style: DECLARATIVE_FRIENDLY
					};
				}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(standardMethodsOnly.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
