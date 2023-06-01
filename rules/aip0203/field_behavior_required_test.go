// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

const ()

func TestFieldBehaviorRequired_SingleFile_SingleMessage(t *testing.T) {
	testCases := []struct {
		name     string
		Fields   string
		problems testutils.Problems
	}{
		{
			"ValidImmutable",
			"int32 page_count = 1 [(google.api.field_behavior) = IMMUTABLE];",
			nil,
		},
		{
			"ValidRequired",
			"int32 page_count = 1 [(google.api.field_behavior) = REQUIRED];",
			nil,
		},
		{
			"ValidOptional",
			"int32 page_count = 1 [(google.api.field_behavior) = OPTIONAL];",
			nil,
		},
		{
			"ValidOutputOnly",
			"int32 page_count = 1 [(google.api.field_behavior) = OUTPUT_ONLY];",
			nil,
		},
		{
			"ValidOptionalImmutable",
			`int32 page_count = 1 [
				(google.api.field_behavior) = OUTPUT_ONLY,
				(google.api.field_behavior) = OPTIONAL
			];`,
			nil,
		},
		{
			"InvalidEmpty",
			"int32 page_count = 1;",
			testutils.Problems{{Message: "annotation must be set"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
					}
				}

				message UpdateBookRequest {
					{{.Fields}}

					Book book = 2 [(google.api.field_behavior) = REQUIRED];
				}

				message UpdateBookResponse {
					// verifies that no error was raised on lack
					// of field behavior in the response.
					string name = 1;
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
					};

					string name = 1;
				}
			`, tc)

			field := f.GetMessageTypes()[0].GetFields()[0]

			if diff := tc.problems.SetDescriptor(field).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFieldBehaviorRequired_Resource_SingleFile(t *testing.T) {
	testCases := []struct {
		name          string
		FieldBehavior string
		problems      testutils.Problems
	}{
		{
			name:          "valid",
			FieldBehavior: "[(google.api.field_behavior) = OUTPUT_ONLY]",
			problems:      nil,
		},
		{
			name:          "invalid",
			FieldBehavior: "",
			problems:      testutils.Problems{{Message: "annotation must be set"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";

				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
					}
				}

				message GetBookRequest {
					string name = 1 [(google.api.field_behavior) = REQUIRED];
				}

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
					};

					string name = 1 {{.FieldBehavior}};
				}
			`, tc)

			field := f.GetMessageTypes()[1].GetFields()[0]

			if diff := tc.problems.SetDescriptor(field).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}

}

func TestFieldBehaviorRequired_NestedMessages_SingleFile(t *testing.T) {
	testCases := []struct {
		name     string
		Fields   string
		problems testutils.Problems
	}{
		{
			"ValidAnnotatedAndChildAnnotated",
			"Annotated annotated = 1 [(google.api.field_behavior) = REQUIRED];",
			nil,
		},
		{
			"InvalidChildNotAnnotated",
			"NonAnnotated non_annotated = 1 [(google.api.field_behavior) = REQUIRED];",
			testutils.Problems{{Message: "must be set"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
					}
				}

				message NonAnnotated {
					string nested = 1;
				}

				message Annotated {
					string nested = 1 [(google.api.field_behavior) = REQUIRED];
				}

				message UpdateBookRequest {
					{{.Fields}}
				}

				message UpdateBookResponse {
					// verifies that no error was raised on lack
					// of field behavior in the response.
					string name = 1;
				}
			`, tc)

			it := f.GetServices()[0].GetMethods()[0].GetInputType()
			nestedField := it.GetFields()[0].GetMessageType().GetFields()[0]

			if diff := tc.problems.SetDescriptor(nestedField).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFieldBehaviorRequired_NestedMessages_MultipleFile(t *testing.T) {
	testCases := []struct {
		name             string
		MessageType      string
		MessageFieldName string
		problems         testutils.Problems
	}{
		{
			"ValidAnnotatedAndChildAnnotated",
			"Annotated",
			"annotated",
			nil,
		},
		{
			"InvalidChildNotAnnotated",
			"NonAnnotated",
			"non_annotated",
			testutils.Problems{{Message: "must be set"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f1 := `
				import "google/api/field_behavior.proto";
				import "resource.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
					}
				}

				message UpdateBookRequest {
					{{.MessageType}} {{.MessageFieldName}} = 1 [(google.api.field_behavior) = REQUIRED];
				}

				message UpdateBookResponse {
					// verifies that no error was raised on lack
					// of field behavior in the response.
					string name = 1;
				}
			`

			f2 := `
				import "google/api/field_behavior.proto";

				message NonAnnotated {
					string nested = 1;
				}

				message Annotated {
					string nested = 1 [(google.api.field_behavior) = REQUIRED];
				}
			`

			srcs := map[string]string{
				"service.proto":  f1,
				"resource.proto": f2,
			}

			ds := testutils.ParseProto3Tmpls(t, srcs, tc)
			f := ds["service.proto"]
			it := f.GetServices()[0].GetMethods()[0].GetInputType()
			fd := it.GetFields()[0].GetMessageType().GetFields()[0]

			if diff := tc.problems.SetDescriptor(fd).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}

			want := "resource.proto"
			if got := fd.GetFile().GetName(); got != want {
				t.Fatalf("got file name %q for location of field but wanted %q", got, want)
			}
		})
	}
}
