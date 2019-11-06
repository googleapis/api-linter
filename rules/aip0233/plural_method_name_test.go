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

package aip0233

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestMethodPluralResourceName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		methodName string
		uriSuffix  string
		problems   testutils.Problems
	}{
		{
			testName:   "Valid-BatchCreateBooks",
			methodName: "BatchCreateBooks",
			uriSuffix:  "books:batchCreate",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Valid-BatchCreateMen",
			methodName: "BatchCreateMen",
			uriSuffix:  "men:batchCreate",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Invalid-SingularBus",
			methodName: "BatchCreateBus",
			uriSuffix:  "bus:batchCreate",
			problems:   testutils.Problems{{Message: `The resource part in method name "BatchCreateBus" shouldn't be "Bus", but should be its plural form "Buses"`}},
		},
		{
			testName:   "Invalid-SingularCorpPerson",
			methodName: "BatchCreateCorpPerson",
			uriSuffix:  "corpPerson:batchCreate",
			problems:   testutils.Problems{{Message: `The resource part in method name "BatchCreateCorpPerson" shouldn't be "CorpPerson", but should be its plural form "CorpPeople"`}},
		},
		{
			testName:   "Invalid-Irrelevant",
			methodName: "AcquireBook",
			uriSuffix:  "book",
			problems:   testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `import "google/api/annotations.proto";

service BookService {
	rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/{{.UriSuffix}}"
			body: "*"
		};
	}
}

message {{.MethodName}}Request {}

message {{.MethodName}}Response{}
`
			file := testutils.ParseProto3Tmpl(t, template,
				struct {
					MethodName string
					UriSuffix  string
				}{test.methodName, test.uriSuffix})

			m := file.GetServices()[0].GetMethods()[0]

			problems := pluralMethodName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
