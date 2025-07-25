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

package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestMethodSignature(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		Signature  string
		problems   testutils.Problems
	}{
		{"Valid", "GetBook", `option (google.api.method_signature) = "name";`, testutils.Problems{}},
		{"Missing", "GetBook", "", testutils.Problems{{Message: `(google.api.method_signature) = "name"`}}},
		{
			"Wrong",
			"GetBook",
			`option (google.api.method_signature) = "book";`,
			testutils.Problems{{Suggestion: `option (google.api.method_signature) = "name";`}},
		},
		{"Irrelevant", "ReadBook", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/client.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book) {
						{{.Signature}}
					}
				}
				message {{.MethodName}}Request {}
				message Book {}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(methodSignature.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
