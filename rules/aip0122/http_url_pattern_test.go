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

package aip0122

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpUriField(t *testing.T) {
	tests := []struct {
		testName string
		src      string
		problems testutils.Problems
	}{
		{
			testName: "Valid",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc CreateInstanceGroup(CreateInstanceGroupsRequest) returns (InstanceGroup) {
		option (google.api.http) = {
			post:"/v1/{parent=publishers/*/instanceGroups/*}"
			body:"instanceGroup"
		};
	}
}

message CreateInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: nil,
		},
		{
			testName: "InValid-Get",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc GetInstanceGroup(GetInstanceGroupsRequest) returns (InstanceGroup) {
	 option (google.api.http) = {
		 get: "/v1/{name=publishers/*/instance_groups/*}"
	 };
	}
}

message GetInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
		{
			testName: "InValid-Create",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc CreateInstanceGroup(CreateInstanceGroupsRequest) returns (InstanceGroup) {
	 option (google.api.http) = {
		 post: "/v1/{parent=publishers/*/instance_groups/*}"
		 body: "instanceGroup"
	 };
	}
}

message CreateInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
		{
			testName: "InValid-Delete",
			//addionalImport:`import "google/protobuf/empty.proto";`,
			src: `import "google/api/annotations.proto";
import "google/protobuf/empty.proto";

service InstanceGroupService {
	rpc DeleteInstanceGroup(DeleteInstanceGroupsRequest) returns (google.protobuf.Empty) {
	 option (google.api.http) = {
		 delete: "/v1/{name=publishers/*/instance_groups/*}"
	 };
	}
}

message DeleteInstanceGroupsRequest {}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
		{
			testName: "InValid-Update_Patch",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc UpdateInstanceGroup(UpdateInstanceGroupsRequest) returns (InstanceGroup) {
	 option (google.api.http) = {
		 patch: "/v1/{book.name=publishers/*/instance_groups/*}"
		 body: "instanceGroup"
	 };
	}
}

message UpdateInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
		{
			testName: "InValid-Update_Put",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc UpdateInstanceGroup(UpdateInstanceGroupsRequest) returns (InstanceGroup) {
	 option (google.api.http) = {
		 put: "/v1/{book.name=publishers/*/instance_groups/*}"
		 body: "book"
	 };
	}
}

message UpdateInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
		{

			testName: "InValid-Custom",
			src: `import "google/api/annotations.proto";

service InstanceGroupService {
	rpc MoveInstanceGroup(MoveInstanceGroupsRequest) returns (InstanceGroup) {
	 option (google.api.http) = {
			custom: { path: "/v1/{book.name=publishers/*/instance_groups/*}" }
	 };
	}
}

message MoveInstanceGroupsRequest {}

message InstanceGroup{}
`,
			problems: testutils.Problems{{
				Message: "HTTP URL should use camel case, but not snake case.",
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)

			f := file.GetServices()[0].GetMethods()[0]
			problems := httpURLPattern.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
