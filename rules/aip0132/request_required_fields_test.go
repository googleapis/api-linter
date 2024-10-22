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

package aip0132

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
			"ValidOptionalPageSize",
			"int32 page_size = 2 [(google.api.field_behavior) = OPTIONAL];",
			"page_size",
			nil,
		},
		{
			"InvalidRequiredPageSize",
			"int32 page_size = 2 [(google.api.field_behavior) = REQUIRED];",
			"page_size",
			testutils.Problems{
				{Message: `List RPCs must only require fields explicitly described in AIPs, not "page_size"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(google.api.field_behavior) = REQUIRED];",
			"create_iam",
			testutils.Problems{
				{Message: `List RPCs must only require fields explicitly described in AIPs, not "create_iam"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";

				service Library {
					rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
						option (google.api.http) = {
							get: "/v1/{parent=publishers/*}/books"
						};
					}
				}

				message ListBooksRequest {
					// The parent, which owns this collection of books.
					// Format: publishers/{publisher}
					string parent = 1 [
					    (google.api.field_behavior) = REQUIRED,
					    (google.api.resource_reference) = {
					  		child_type: "library.googleapis.com/Book"
					    }];

					{{.Fields}}
				}

				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			var dbr desc.Descriptor = f.FindMessage("ListBooksRequest")
			if test.problematicFieldName != "" {
				dbr = f.FindMessage("ListBooksRequest").FindFieldByName(test.problematicFieldName)
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
