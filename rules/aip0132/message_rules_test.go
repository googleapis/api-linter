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

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestStandardFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		nameFieldName string
		nameFieldType *builder.FieldType
		problemCount  int
		errPrefix     string
	}{
		{"Valid", "ListBooksRequest", "parent", builder.FieldTypeString(), 0, "False positive"},
		{"InvalidName", "ListBooksRequest", "publisher", builder.FieldTypeString(), 1, "False negative"},
		{"InvalidType", "ListBooksRequest", "parent", builder.FieldTypeBytes(), 1, "False negative"},
		{"Irrelevant", "EnumerateBooksRequest", "id", builder.FieldTypeString(), 0, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField(test.nameFieldName, test.nameFieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := standardFields.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, standardFields.Name, problems)
			}
		})
	}
}

func TestUnknownFields(t *testing.T) {
	// Get the correct message type for google.protobuf.FieldMask.
	fieldMask, err := desc.LoadMessageDescriptorForMessage(&field_mask.FieldMask{})
	if err != nil {
		t.Fatalf("Unable to load the field mask message.")
	}

	// Set up the testing permutations.
	tests := []struct {
		testName     string
		messageName  string
		fieldName    string
		fieldType    *builder.FieldType
		problemCount int
		errPrefix    string
	}{
		{"PageSize", "ListBooksRequest", "page_size", builder.FieldTypeInt32(), 0, "False positive"},
		{"PageToken", "ListBooksRequest", "page_token", builder.FieldTypeString(), 0, "False positive"},
		{"Filter", "ListBooksRequest", "filter", builder.FieldTypeString(), 0, "False positive"},
		{"FilterInvalid", "ListBooksRequest", "filter", builder.FieldTypeBytes(), 1, "False negative"},
		{"OrderBy", "ListBooksRequest", "order_by", builder.FieldTypeString(), 0, "False positive"},
		{"OrderByInvalid", "ListBooksRequest", "order_by", builder.FieldTypeBytes(), 1, "False negative"},
		{"GroupBy", "ListBooksRequest", "group_by", builder.FieldTypeString(), 0, "False positive"},
		{"GroupByInvalid", "ListBooksRequest", "group_by", builder.FieldTypeBytes(), 1, "False negative"},
		{"ShowDeleted", "ListBooksRequest", "show_deleted", builder.FieldTypeBool(), 0, "False positive"},
		{"ReadMask", "ListBooksRequest", "read_mask", builder.FieldTypeImportedMessage(fieldMask), 0, "False positive"},
		{"View", "ListBooksRequest", "view", builder.FieldTypeEnum(builder.NewEnum("View")), 0, "False positive"},
		{"Invalid", "ListBooksRequest", "application_id", builder.FieldTypeString(), 1, "False negative"},
		{"Irrelevant", "EnumerteBooksRequest", "application_id", builder.FieldTypeString(), 0, "False positive"},
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
				t.Fatalf("Could not build GetBookRequest message.")
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := unknownFields.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, unknownFields.Name, problems)
			}
		})
	}
}
