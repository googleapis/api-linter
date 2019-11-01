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

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestPluralMethodResourceName(t *testing.T) {
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

			problems := pluralMethodResourceName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestInputName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName string
		src      string
		problems testutils.Problems
	}{
		{
			testName: "Valid-BatchCreateBooksRequest",
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
			testName: "Valid-BatchCreateMenRequest",
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
	rpc BatchCreateBuses(BatchCreateBusRequest) returns (BatchCreateBusResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/buses:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateBusRequest {}

message BatchCreateBusResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Create RPCs should have a properly named request message "BatchCreateBusesRequest", but not "BatchCreateBusRequest"`}},
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
			problems: testutils.Problems{{Message: `Batch Create RPCs should have a properly named request message "BatchCreateCorpPeopleRequest", but not "BatchCreateCorpPersonRequest"`}},
		},
		{
			testName: "Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc CreateBook(CreateBookRequest) returns (Book) {
		option (google.api.http) = {
			post: "/v1/{name=publishers/*/books/*}"
			body: "*"
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

			problems := inputName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestOutputName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName string
		src      string
		problems testutils.Problems
	}{
		{
			testName: "Valid-BatchCreateBooksResponse",
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
			testName: "Valid-LongRunning",
			src: `import "google/api/annotations.proto";
import "google/longrunning/operations.proto";

service BookService {
	rpc BatchCreateBooks(BatchCreateBooksRequest) returns (google.longrunning.Operation) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/books:batchCreate"
			body: "*"
		};
		option (google.longrunning.operation_info) = {
      response_type: "BatchCreateBooksResponse"
    };
	}
}

message BatchCreateBooksRequest {}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-LongRunningEmptyResponseType",
			src: `import "google/api/annotations.proto";
import "google/longrunning/operations.proto";

service BookService {
	rpc BatchCreateBooks(BatchCreateBooksRequest) returns (google.longrunning.Operation) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/books:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateBooksRequest {}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-BatchCreateMenResponse",
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
	rpc BatchCreateBuses(BatchCreateBusRequest) returns (BatchCreateBusResponse) {
		option (google.api.http) = {
			post: "/v1/{parent=publishers/*}/buses:batchCreate"
			body: "*"
		};
	}
}

message BatchCreateBusRequest {}

message BatchCreateBusResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Create RPCs should have a properly named response message "BatchCreateBusesResponse", but not "BatchCreateBusResponse"`}},
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
			problems: testutils.Problems{{Message: `Batch Create RPCs should have a properly named response message "BatchCreateCorpPeopleResponse", but not "BatchCreateCorpPersonResponse"`}},
		},
		{
			testName: "Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc CreateBook(CreateBookRequest) returns (Book) {
		option (google.api.http) = {
			post: "/v1/{name=publishers/*/books/*}"
			body: "*"
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

			problems := outputName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpVerb(t *testing.T) {
	// Set up GET and POST HTTP annotations.
	httpGet := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{parent=publishers/*}/books:batchGet",
		},
	}
	httpCreate := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/{parent=publishers/*}/books:batchCreate",
		},
		Body: "*",
	}

	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpRule   *annotations.HttpRule
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", httpCreate, "BatchCreateBooks", nil},
		{"Invalid", httpGet, "BatchCreateBooks", testutils.Problems{{Message: "Batch Create methods must use the HTTP POST verb."}}},
		{"Irrelevant", httpGet, "AcquireBook", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-233 Create method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksResponse"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpVerb.Lint(service.GetFile())
			if diff := test.problems.SetDescriptor(service.GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpUrl(t *testing.T) {
	tests := []struct {
		testName   string
		uri        string
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", "/v1/{parent=publishers/*}/books:batchCreate", "BatchCreateBooks", nil},
		{"InvalidVarName", "/v1/{parent=publishers/*}/books", "BatchCreateBooks", testutils.Problems{{Message: `Batch Create methods URI should be end with ":batchCreate".`}}},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "AcquireBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: test.uri,
				},
				Body: "*",
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-233 Create method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksResponse"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpUrl.Lint(service.GetFile())
			if diff := test.problems.SetDescriptor(service.GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName   string
		body       string
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", "*", "BatchCreateBooks", nil},
		{"Invalid", "", "BatchCreateBooks", testutils.Problems{{Message: `Batch Create methods should use "*" as the HTTP body.`}}},
		{"Irrelevant", "*", "AcquireBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{parent=publishers/*}/books:batchCreate",
				},
				Body: test.body,
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-233 Create method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchCreateBooksResponse"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpBody.Lint(service.GetFile())
			if diff := test.problems.SetDescriptor(service.GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
