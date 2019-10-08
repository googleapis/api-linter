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

package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	fpb "google.golang.org/genproto/protobuf/field_mask"
)

func TestStandardFields(t *testing.T) {
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
		problems    testutils.Problems
	}{
		// We use BigBook instead of Book in order to test correct casing logic
		{"Valid", "UpdateBigBookRequest", "big_book", testutils.Problems{}},
		{"NoResource", "UpdateBigBookRequest", "id", testutils.Problems{{Message: "book"}}},
		{"Irrelevant", "AcquireBigBookRequest", "id", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			bookMsg := builder.NewMessage("BigBook")
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField(test.fieldName, builder.FieldTypeMessage(bookMsg)),
			).AddField(
				builder.NewField("update_mask", builder.FieldTypeImportedMessage(fieldMask)),
			).Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := standardFields.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}

func TestStandardFieldsMissingUpdateMask(t *testing.T) {
	// Create an appropriate message descriptor.
	bookMsg := builder.NewMessage("Book")
	message, err := builder.NewMessage("UpdateBookRequest").AddField(
		builder.NewField("book", builder.FieldTypeMessage(bookMsg)),
	).Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Run the lint rule, and establish that it returns the correct
	// number of problems.
	wantProblems := testutils.Problems{{
		Descriptor: message,
		Message:    "Method UpdateBookRequest has no `update_mask` field",
	}}
	gotProblems := standardFields.Lint(message.GetFile())
	if diff := wantProblems.Diff(gotProblems); diff != "" {
		t.Errorf(diff)
	}
}

func TestStandardFieldsInvalidType(t *testing.T) {
	// Create an appropriate message descriptor.
	parchmentMsg := builder.NewMessage("Parchment")
	message, err := builder.NewMessage("UpdateBookRequest").AddField(
		builder.NewField("book", builder.FieldTypeMessage(parchmentMsg)),
	).Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Run the lint rule, and establish that it returns the correct
	// number of problems.
	wantProblems := testutils.Problems{{
		Descriptor: message.GetFields()[0],
		Message:    "`book` field on Update RPCs should be of type `Book`",
	}}
	gotProblems := standardFields.Lint(message.GetFile())
	if diff := wantProblems.Diff(gotProblems); diff != "" {
		t.Errorf(diff)
	}
}

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
		// Use BigBook instead of Book to test correct casing logic
		{"UpdateMask", "UpdateBigBookRequest", "update_mask",
			builder.FieldTypeImportedMessage(fieldMask), testutils.Problems{}},
		{"Invalid", "UpdateBigBookRequest", "application_id",
			builder.FieldTypeString(), testutils.Problems{{Message: "Unexpected field"}}},
		{"InvalidCasing", "UpdateBigBookRequest", "bigbook",
			builder.FieldTypeString(), testutils.Problems{{Message: "Unexpected field"}}},
		{"Irrelevant", "AcquireBigBookRequest", "application_id",
			builder.FieldTypeString(), testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField("big_book", builder.FieldTypeMessage(builder.NewMessage("BigBook"))),
			).AddField(
				builder.NewField(test.fieldName, test.fieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build UpdateBookRequest message.")
			}

			// Run the lint rule, and establish that it returns the correct problems.
			wantProblems := test.problems.SetDescriptor(message.FindFieldByName(test.fieldName))
			gotProblems := unknownFields.Lint(message.GetFile())
			if diff := wantProblems.Diff(gotProblems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
