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

package aip0134

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
			"ValidOptionalUpdateMask",
			"google.protobuf.FieldMask update_mask = 2 [(google.api.field_behavior) = OPTIONAL];",
			"update_mask",
			nil,
		},
		{
			"ValidRequiredUpdateMask",
			"google.protobuf.FieldMask update_mask = 2 [(google.api.field_behavior) = REQUIRED];",
			"update_mask",
			nil,
		},
		{
			"ValidOptionalValidateOnly",
			"bool validate_only = 3 [(google.api.field_behavior) = OPTIONAL];",
			"validate_only",
			nil,
		},
		{
			"InvalidRequiredValidateOnly",
			"bool validate_only = 3 [(google.api.field_behavior) = REQUIRED];",
			"validate_only",
			testutils.Problems{
				{Message: `Update RPCs must only require fields explicitly described in AIPs, not "validate_only"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(google.api.field_behavior) = REQUIRED];",
			"create_iam",
			testutils.Problems{
				{Message: `Update RPCs must only require fields explicitly described in AIPs, not "create_iam"`},
			},
		},
		{
			"InvalidRequiredUnknownMessageField",
			"Foo foo = 3 [(google.api.field_behavior) = REQUIRED];",
			"foo",
			testutils.Problems{
				{Message: `Update RPCs must only require fields explicitly described in AIPs, not "foo"`},
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
					rpc UpdateBookShelf(UpdateBookShelfRequest) returns (BookShelf) {
						option (google.api.http) = {
							patch: "/v1/{name=publishers/*/bookShelves/*}"
							body: "book"
						};
					}
				}

				message BookShelf {
					option (google.api.resource) = {
						type: "library.googleapis.com/BookShelf"
						pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
					};
					string name = 1;
				}

				message UpdateBookShelfRequest {
					BookShelf book_shelf = 1 [
						(google.api.field_behavior) = REQUIRED
					];
					{{.Fields}}
				}

				message Foo {}
			`, test)
			var dbr desc.Descriptor = f.FindMessage("UpdateBookShelfRequest")
			if test.problematicFieldName != "" {
				dbr = f.FindMessage("UpdateBookShelfRequest").FindFieldByName(test.problematicFieldName)
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
