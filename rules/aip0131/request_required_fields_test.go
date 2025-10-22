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

package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
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
			"ValidOptionalReadMask",
			"google.protobuf.FieldMask read_mask = 2 [(google.api.field_behavior) = OPTIONAL];",
			"read_mask",
			nil,
		},
		{
			"InvalidRequiredReadMask",
			"google.protobuf.FieldMask read_mask = 2 [(google.api.field_behavior) = REQUIRED];",
			"read_mask",
			testutils.Problems{
				{Message: `Get RPCs must only require fields explicitly described in AIPs, not "read_mask"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(google.api.field_behavior) = REQUIRED];",
			"create_iam",
			testutils.Problems{
				{Message: `Get RPCs must only require fields explicitly described in AIPs, not "create_iam"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";

				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "/v1/{name=publishers/*/books/*}"
						};
					}
				}

				message GetBookRequest {
					// The name of the book to retrieve.
					// Format: publishers/{publisher}/books/{book}
					string name = 1 [
					    (google.api.field_behavior) = REQUIRED,
						(google.api.resource_reference) = {
							type: "library.googleapis.com/Book"
						}
					];
					{{.Fields}}
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			var dbr protoreflect.Descriptor = f.Messages().Get(0)
			if test.problematicFieldName != "" {
				dbr = f.Messages().Get(0).Fields().ByName(protoreflect.Name(test.problematicFieldName))
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
