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

package aip0234

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpMethod(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpVerb   string
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", "post", "BatchUpdateBooks", nil},
		{"Invalid", "get", "BatchUpdateBooks", testutils.Problems{{Message: "HTTP POST verb"}}},
		{"Irrelevant", "get", "AcquireBook", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `import "google/api/annotations.proto";
service BookService {
	rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
		option (google.api.http) = {
			{{.HttpVerb}}: "/v1/{parent=publishers/*}/books:batchUpdate"
			body: "*"
		};
	}
}
message {{.MethodName}}Request{}
message {{.MethodName}}Response{}
`
			file := testutils.ParseProto3Tmpl(t, template,
				struct {
					MethodName string
					HttpVerb   string
				}{test.methodName, test.httpVerb})

			// Run the method, ensure we get what we expect.
			problems := httpMethod.Lint(file)
			if diff := test.problems.SetDescriptor(file.GetServices()[0].GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
