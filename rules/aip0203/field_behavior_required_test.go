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

func TestFieldBehaviorRequired(t *testing.T) {
	for _, test := range []struct {
		name     string
		Fields   string
		problems testutils.Problems
	}{
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
			"InvalidImmutable",
			"int32 page_count = 1 [(google.api.field_behavior) = IMMUTABLE];",
			testutils.Problems{{Message: "must have at least one"}},
		},
		{
			"InvalidEmpty",
			"int32 page_count = 1;",
			testutils.Problems{{Message: "annotation must be set"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
					}
				}

				message UpdateBookRequest {
					{{.Fields}}
				}

				message UpdateBookResponse {
					// verifies that no error was raised on lack
					// of field behavior in the response.
					string name = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFieldBehaviorRequiredNested(t *testing.T) {
	for _, test := range []struct {
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
	} {
		t.Run(test.name, func(t *testing.T) {
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
			`, test)
			it := f.GetServices()[0].GetMethods()[0].GetInputType()
			nestedField := it.GetFields()[0].GetMessageType().GetFields()[0]
			if diff := test.problems.SetDescriptor(nestedField).Diff(fieldBehaviorRequired.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
