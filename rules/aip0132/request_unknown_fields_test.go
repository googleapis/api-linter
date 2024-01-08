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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	fpb "google.golang.org/genproto/protobuf/field_mask"
)

func TestUnknownFields(t *testing.T) {
	// Get the correct message type for google.protobuf.FieldMask.
	fieldMask, err := desc.LoadMessageDescriptorForMessage(&fpb.FieldMask{})
	if err != nil {
		t.Fatalf("Unable to load the field mask message.")
	}

	// Set up the testing permutations.
	tests := []struct {
		testName    string
		messageName string
		fieldName   string
		fieldType   *builder.FieldType
		problems    testutils.Problems
	}{
		{"PageSize", "ListBooksRequest", "page_size", builder.FieldTypeInt32(), testutils.Problems{}},
		{"PageToken", "ListBooksRequest", "page_token", builder.FieldTypeString(), testutils.Problems{}},
		{"Skip", "ListBooksRequest", "skip", builder.FieldTypeInt32(), testutils.Problems{}},
		{"Filter", "ListBooksRequest", "filter", builder.FieldTypeString(), testutils.Problems{}},
		{"OrderBy", "ListBooksRequest", "order_by", builder.FieldTypeString(), testutils.Problems{}},
		{"ShowDeleted", "ListBooksRequest", "show_deleted", builder.FieldTypeBool(), testutils.Problems{}},
		{"RequestId", "ListBooksRequest", "request_id", builder.FieldTypeImportedMessage(fieldMask), testutils.Problems{}},
		{"ReadMask", "ListBooksRequest", "read_mask", builder.FieldTypeImportedMessage(fieldMask), testutils.Problems{}},
		{"View", "ListBooksRequest", "view", builder.FieldTypeEnum(builder.NewEnum("View").AddValue(builder.NewEnumValue("BASIC"))), testutils.Problems{}},
		{"Invalid", "ListBooksRequest", "application_id", builder.FieldTypeString(), testutils.Problems{{Message: "explicitly described"}}},
		{"Irrelevant", "EnumerteBooksRequest", "application_id", builder.FieldTypeString(), testutils.Problems{}},
		{"IrrelevantAIP162", "ListBookRevisionsRequest", "name", builder.FieldTypeString(), testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField("parent", builder.FieldTypeString()),
			).AddField(
				builder.NewField(test.fieldName, test.fieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build GetBookRequest message: %s", err)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			problems := unknownFields.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(message.GetFields()[1]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
