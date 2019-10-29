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

package aip0231

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
			testName: "Valid-BatchGetBooks",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/books:batchGet"
		};
	}
}

message BatchGetBooksRequest {}

message BatchGetBooksResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-BatchGetMen",
			src: `import "google/api/annotations.proto";

service ManService {
	rpc BatchGetMen(BatchGetMenRequest) returns (BatchGetMenResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/men:batchGet"
		};
	}
}

message BatchGetMenRequest {}

message BatchGetMenResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-SingularBus",
			src: `import "google/api/annotations.proto";

service BusService {
	rpc BatchGetBus(BatchGetBusRequest) returns (BatchGetBusResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/buses:batchGet"
		};
	}
}

message BatchGetBusRequest {}

message BatchGetBusResponse{}
`,
			problems: testutils.Problems{{Message: `The resource part in method name "BatchGetBus" shouldn't "Bus", but should be its plural form "Buses"`}},
		},
		{
			testName: "Invalid-SingularCorpPerson",
			src: `import "google/api/annotations.proto";

service CorpPersonService {
	rpc BatchGetCorpPerson(BatchGetCorpPersonRequest) returns (BatchGetCorpPersonResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/corpPerson:batchGet"
		};
	}
}

message BatchGetCorpPersonRequest {}

message BatchGetCorpPersonResponse{}
`,
			problems: testutils.Problems{{Message: `The resource part in method name "BatchGetCorpPerson" shouldn't "CorpPerson", but should be its plural form "CorpPeople"`}},
		},
		{
			testName: "Invalid-Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc GetBook(GetBookRequest) returns (Book) {
		option (google.api.http) = {
			get: "/v1/{name=publishers/*/books/*}"
		};
	}
}

message GetBookRequest {}

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
			testName: "Valid-BatchGetBooksRequest",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/books:batchGet"
		};
	}
}

message BatchGetBooksRequest {}

message BatchGetBooksResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-BatchGetMenRequest",
			src: `import "google/api/annotations.proto";

service ManService {
	rpc BatchGetMen(BatchGetMenRequest) returns (BatchGetMenResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/men:batchGet"
		};
	}
}

message BatchGetMenRequest {}

message BatchGetMenResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-SingularBus",
			src: `import "google/api/annotations.proto";

service BusService {
	rpc BatchGetBuses(BatchGetBusRequest) returns (BatchGetBusResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/buses:batchGet"
		};
	}
}

message BatchGetBusRequest {}

message BatchGetBusResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Get RPCs should have a properly named request message "BatchGetBusesRequest", but not "BatchGetBusRequest"`}},
		},
		{
			testName: "Invalid-SingularCorpPerson",
			src: `import "google/api/annotations.proto";

service CorpPersonService {
	rpc BatchGetCorpPerson(BatchGetCorpPersonRequest) returns (BatchGetCorpPersonResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/corpPerson:batchGet"
		};
	}
}

message BatchGetCorpPersonRequest {}

message BatchGetCorpPersonResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Get RPCs should have a properly named request message "BatchGetCorpPeopleRequest", but not "BatchGetCorpPersonRequest"`}},
		},
		{
			testName: "Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc GetBook(GetBookRequest) returns (Book) {
		option (google.api.http) = {
			get: "/v1/{name=publishers/*/books/*}"
		};
	}
}

message GetBookRequest {}

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
			testName: "Valid-BatchGetBooksResponse",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/books:batchGet"
		};
	}
}

message BatchGetBooksRequest {}

message BatchGetBooksResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-BatchGetBooksResponse",
			src: `import "google/api/annotations.proto";

service ManService {
	rpc BatchGetMen(BatchGetMenRequest) returns (BatchGetMenResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/men:batchGet"
		};
	}
}

message BatchGetMenRequest {}

message BatchGetMenResponse{}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-SingularBus",
			src: `import "google/api/annotations.proto";

service BusService {
	rpc BatchGetBuses(BatchGetBusRequest) returns (BatchGetBusResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/buses:batchGet"
		};
	}
}

message BatchGetBusRequest {}

message BatchGetBusResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Get RPCs should have a properly named response message "BatchGetBusesResponse", but not "BatchGetBusResponse"`}},
		},
		{
			testName: "Invalid-SingularCorpPerson",
			src: `import "google/api/annotations.proto";

service CorpPersonService {
	rpc BatchGetCorpPerson(BatchGetCorpPersonRequest) returns (BatchGetCorpPersonResponse) {
		option (google.api.http) = {
			get: "/v1/{parent=publishers/*}/corpPerson:batchGet"
		};
	}
}

message BatchGetCorpPersonRequest {}

message BatchGetCorpPersonResponse{}
`,
			problems: testutils.Problems{{Message: `Batch Get RPCs should have a properly named response message "BatchGetCorpPeopleResponse", but not "BatchGetCorpPersonResponse"`}},
		},
		{
			testName: "Irrelevant",
			src: `import "google/api/annotations.proto";

service BookService {
	rpc GetBook(GetBookRequest) returns (Book) {
		option (google.api.http) = {
			get: "/v1/{name=publishers/*/books/*}"
		};
	}
}

message GetBookRequest {}

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
	httpPost := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/{name=publishers/*/books/*}",
		},
	}

	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpRule   *annotations.HttpRule
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", httpGet, "BatchGetBooks", nil},
		{"Invalid", httpPost, "BatchGetBooks", testutils.Problems{{Message: "Batch Get methods must use the HTTP GET verb."}}},
		{"Irrelevant", httpPost, "AcquireBook", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-231 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksResponse"), false),
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
		{"Valid", "/v1/{parent=publishers/*}/books:batchGet", "BatchGetBooks", nil},
		{"InvalidVarName", "/v1/{parent=publishers/*}/books", "BatchGetBooks", testutils.Problems{{Message: `Get methods URI should be end with ":batchGet".`}}},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "AcquireBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{
					Get: test.uri,
				},
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-231 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksResponse"), false),
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
		{"Valid", "", "BatchGetBooks", nil},
		{"Invalid", "*", "BatchGetBooks", testutils.Problems{{Message: "Batch Get methods should not have an HTTP body."}}},
		{"Irrelevant", "*", "AcquireBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{
					Get: "/v1/{name=publishers/*/books/*}",
				},
				Body: test.body,
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-231 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksResponse"), false),
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
