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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRequiredFieldTests(t *testing.T) {
	for _, test := range []struct {
		name                 string
		Fields               string
		problematicFieldName string
		Singular             string
		problems             testutils.Problems
	}{
		{
			"ValidNoExtraFields",
			"",
			"",
			"",
			nil,
		},
		{
			"ValidWithSingularNoExtraFields",
			"",
			"",
			"bookShelf",
			nil,
		},
		{
			"ValidWithSingularAndIdField",
			"string book_shelf_id = 3 [(google.api.field_behavior) = OPTIONAL];",
			"",
			"bookShelf",
			nil,
		},
		{
			"ValidOptionalValidateOnly",
			"string validate_only = 3 [(google.api.field_behavior) = OPTIONAL];",
			"validate_only",
			"",
			nil,
		},
		{
			"InvalidRequiredValidateOnly",
			"bool validate_only = 3 [(google.api.field_behavior) = REQUIRED];",
			"validate_only",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AIPs, not "validate_only"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(google.api.field_behavior) = REQUIRED];",
			"create_iam",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AIPs, not "create_iam"`},
			},
		},
		{
			"InvalidRequiredUnknownMessageField",
			"Foo foo = 3 [(google.api.field_behavior) = REQUIRED];",
			"foo",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AIPs, not "foo"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";

				service Library {
					rpc CreateBookShelf(CreateBookShelfRequest) returns (BookShelf) {
						option (google.api.http) = {
							delete: "/v1/{name=publishers/*/bookShelves/*}"
						};
					}
				}

				message BookShelf {
					option (google.api.resource) = {
						type: "library.googleapis.com/BookShelf"
						pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
						singular: "{{.Singular}}"
					};
					string name = 1;
				}

				message Foo {}

				message CreateBookShelfRequest {
					string parent = 1 [
						(google.api.field_behavior) = REQUIRED
					];
					BookShelf book_shelf = 2 [
						(google.api.field_behavior) = REQUIRED
					];
					{{.Fields}}
				}
			`, test)
			var dbr protoreflect.Descriptor = f.Messages().Get(2)
			if test.problematicFieldName != "" {
				dbr = f.Messages().Get(2).Fields().ByName(protoreflect.Name(test.problematicFieldName))
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
