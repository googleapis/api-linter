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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName             string
		MessageName          string
		FieldName            string
		FieldType            string
		problems             testutils.Problems
		problematicFieldName string
	}{
		{
			"Parent",
			"CreateBookRequest",
			"parent",
			"string",
			testutils.Problems{},
			"",
		},
		{
			"ValidateOnly",
			"CreateBookRequest",
			"validate_only",
			"bool",
			testutils.Problems{},
			"",
		},
		{
			"ResourceRelatedField",
			"CreateBookRequest",
			"book_id",
			"string",
			testutils.Problems{},
			"",
		},
		{
			"ResourceRelatedField",
			"CreateBookStoreRequest",
			"book_store_id",
			"string",
			testutils.Problems{},
			"",
		},
		{
			"RequestIdField",
			"CreateBookRequest",
			"request_id",
			"string",
			testutils.Problems{},
			"",
		},
		{
			"Invalid",
			"CreateBookRequest",
			"name",
			"string",
			testutils.Problems{{Message: "Create RPCs must only contain fields explicitly described in AIPs, not \"name\"."}},
			"name",
		},
		{
			"InvalidResourceRelatedField",
			"CreateBookStoreRequest",
			"book_id",
			"string",
			testutils.Problems{{Message: "Create RPCs must only contain fields explicitly described in AIPs, not \"book_id\"."}},
			"book_id",
		},
		{
			"Irrelevant",
			"GetBookRequest",
			"name",
			"string",
			testutils.Problems{},
			"",
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)

			problems := unknownFields.Lint(f)
			var d protoreflect.Descriptor = f.Messages().Get(0)
			if test.problematicFieldName != "" {
				d = f.Messages().Get(0).Fields().Get(0)
			}

			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}