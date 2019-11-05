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
		testName string
		src      string
		problems testutils.Problems
	}{
		{
			testName: "Valid-BatchCreateBooks",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/books:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateBooksRequest {}

message BatchCreateBooksResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-BatchCreateMen",
			src: `import "google/api/annotations.proto";

service ManService {
	rpc BatchCreateMen(BatchCreateMenRequest) returns (BatchCreateMenResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/men:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateMenRequest {}

message BatchCreateMenResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-SingularBus",
			src: `import "google/api/annotations.proto";

service BusService {
	rpc BatchCreateBus(BatchCreateBusRequest) returns (BatchCreateBusResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/buses:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateBusRequest {}

message BatchCreateBusResponse{}
`,
			problems: testutils.Problems{{Message: `The resource part in method name "BatchCreateBus" shouldn't be "Bus", but should be its plural form "Buses"`}},
		},
		{
			testName: "Invalid-SingularCorpPerson",
			src: `import "google/api/annotations.proto";

service CorpPersonService {
	rpc BatchCreateCorpPerson(BatchCreateCorpPersonRequest) returns (BatchCreateCorpPersonResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/corpPerson:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateCorpPersonRequest {}

message BatchCreateCorpPersonResponse{}
`,
			problems: testutils.Problems{{Message: `The resource part in method name "BatchCreateCorpPerson" shouldn't be "CorpPerson", but should be its plural form "CorpPeople"`}},
		},
		{
			testName: "Invalid-Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc CreateBook(CreateBookRequest) returns (Book) {
		option (google.api.http) = {
			get: "/v1/{name=publishers/*/books/*}"
		};
	}
}

message CreateBookRequest {}

message Book{}
`,
			problems: testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)

			m := file.GetServices()[0].GetMethods()[0]

			problems := pluralMethodName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
