// Copyright 2023 Google LLC
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

package aip0135

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRequiredFieldTests(t *testing.T) {
	for _, test := range []struct {
		name                 string
		Fields               string
		problematicFieldName string
		problems             testutils.Problems
	}{
		{
			"ValidNoExtraFields",
			"",
			"",
			nil,
		},
		{
			"ValidOptionalAllowMissing",
			"bool allow_missing = 2 [(google.api.field_behavior) = OPTIONAL];",
			"allow_missing",
			nil,
		},
		{
			"InvalidRequiredAllowMissing",
			"bool allow_missing = 2 [(google.api.field_behavior) = REQUIRED];",
			"allow_missing",
			testutils.Problems{
				{Message: `Delete RPCs must only require fields explicitly described in AIPs, not "allow_missing"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 2 [(google.api.field_behavior) = REQUIRED];",
			"create_iam",
			testutils.Problems{
				{Message: `Delete RPCs must only require fields explicitly described in AIPs, not "create_iam"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";
			    import "google/longrunning/operations.proto";

				service Library {
					rpc DeleteBook(DeleteBookRequest) returns (google.longrunning.Operation) {
						option (google.api.http) = {
						    delete: "/v1/{name=publishers/*/books/*}"
						};
						option (google.longrunning.operation_info) = {
							response_type: "google.protobuf.Empty"
							metadata_type: "OperationMetadata"
						};
					}
				}

				message DeleteBookRequest {
					string name = 1 [
						(google.api.field_behavior) = REQUIRED
					];
					{{.Fields}}
				}
			`, test)
			var dbr desc.Descriptor = f.FindMessage("DeleteBookRequest")
			if test.problematicFieldName != "" {
				dbr = f.FindMessage("DeleteBookRequest").FindFieldByName(test.problematicFieldName)
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
