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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRequestRequiredFieldsType(t *testing.T) {
	for _, test := range []struct {
		name                 string
		Fields               string
		problematicFieldName string
		Singular             string
		ReturnType           string
		problems             testutils.Problems
	}{
		{
			name: "ValidStandardFields",
			Fields: `string parent = 1;
			 BookShelf book_shelf = 2;
			 string book_shelf_id = 3;`,
			problematicFieldName: "",
			Singular:             "bookShelf",
			ReturnType:           "BookShelf",
			problems:             nil,
		},
		{
			name: "ValidRpcNameInference",
			Fields: `string parent = 1;
			 BookShelf book_shelf = 2;
			 string book_shelf_id = 3;`,
			problematicFieldName: "",
			Singular:             "bookShelf",
			ReturnType:           "CreateBookShelfResponse",
			problems:             nil,
		},
		{
			name: "InvalidParentType",
			Fields: `int32 parent = 1;
			 BookShelf book_shelf = 2;
			 string book_shelf_id = 3;`,
			problematicFieldName: "parent",
			Singular:             "bookShelf",
			ReturnType:           "BookShelf",
			problems: testutils.Problems{
				{Message: `The field "parent" must be of type string.`},
			},
		},
		{
			name: "InvalidResourceIdType",
			Fields: `string parent = 1;
			 BookShelf book_shelf = 2;
			 int32 book_shelf_id = 3;`,
			problematicFieldName: "book_shelf_id",
			Singular:             "bookShelf",
			ReturnType:           "BookShelf",
			problems: testutils.Problems{
				{Message: `The resource ID field "book_shelf_id" must be of type string.`},
			},
		},
		{
			name: "InvalidResourceType",
			Fields: `string parent = 1;
			 string book_shelf = 2;
			 string book_shelf_id = 3;`,
			problematicFieldName: "book_shelf",
			Singular:             "bookShelf",
			ReturnType:           "BookShelf",
			problems: testutils.Problems{
				{Message: `The resource field "book_shelf" must be of type message.`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";

				service Library {
					rpc CreateBookShelf(CreateBookShelfRequest) returns ({{.ReturnType}}) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/bookShelves",
							body: "book_shelf"
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

				message CreateBookShelfRequest {
					{{.Fields}}
				}

				message CreateBookShelfResponse {}
			`, test)

			reqMsg := f.Messages().ByName("CreateBookShelfRequest")
			var dbr protoreflect.Descriptor = reqMsg
			if test.problematicFieldName != "" {
				dbr = reqMsg.Fields().ByName(protoreflect.Name(test.problematicFieldName))
			}

			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFieldsTypes.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
