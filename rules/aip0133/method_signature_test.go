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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestMethodSignature(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		Signature  string
		IDField    string
		problems   testutils.Problems
	}{
		{"ValidNoID", "CreateBook", `option (google.api.method_signature) = "parent,book";`, "", testutils.Problems{}},
		{"ValidID", "CreateBook", `option (google.api.method_signature) = "parent,book,book_id";`, "string book_id = 3;", testutils.Problems{}},
		{"MissingNoID", "CreateBook", "", "", testutils.Problems{{Message: `(google.api.method_signature) = "parent,book"`}}},
		{
			"MissingID",
			"CreateBook",
			"",
			"string book_id = 3;",
			testutils.Problems{{Message: `(google.api.method_signature) = "parent,book,book_id"`}},
		},
		{
			"Wrong",
			"CreateBook",
			`option (google.api.method_signature) = "publisher,book";`,
			"",
			testutils.Problems{{Suggestion: `option (google.api.method_signature) = "parent,book";`}},
		},
		{
			"WrongIDMissing",
			"CreateBook",
			`option (google.api.method_signature) = "parent,book";`,
			"string book_id = 3;",
			testutils.Problems{{Suggestion: `option (google.api.method_signature) = "parent,book,book_id";`}},
		},
		{
			"WrongIDPresent",
			"CreateBook",
			`option (google.api.method_signature) = "parent,book,book_id";`,
			"",
			testutils.Problems{{Suggestion: `option (google.api.method_signature) = "parent,book";`}},
		},
		{"Irrelevant", "WriteBook", "", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/client.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book) {
						{{.Signature}}
					}
				}
				message {{.MethodName}}Request {
					string parent = 1;
					Book book = 2;
					{{.IDField}}
				}
				message Book {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(methodSignature.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
