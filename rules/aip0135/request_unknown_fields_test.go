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

package aip0135

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		messageName string
		fieldName   string
		fieldType   *builder.FieldType
		problems    testutils.Problems
	}{
		{"Force", "DeleteBookRequest", "force", builder.FieldTypeBool(), testutils.Problems{}},
		{"Etag", "DeleteBookRequest", "etag", builder.FieldTypeString(), testutils.Problems{}},
		{"AllowMissing", "DeleteBookRequest", "allow_missing", builder.FieldTypeBool(), testutils.Problems{}},
		{"RequestId", "DeleteBookRequest", "request_id", builder.FieldTypeString(), testutils.Problems{}},
		{"ValidateOnly", "DeleteBookRequest", "validate_only", builder.FieldTypeBool(), testutils.Problems{}},
		{"Invalid", "DeleteBookRequest", "application_id", builder.FieldTypeString(), testutils.Problems{{
			Message: "Unexpected field",
		}}},
		{"Irrelevant", "RemoveBookRequest", "application_id", builder.FieldTypeString(), testutils.Problems{}},
		{"ValidView", "DeleteBookRequest", "view", builder.FieldTypeEnum(
			builder.NewEnum("BookView").
				AddValue(builder.NewEnumValue("BOOK_VIEW_UNSPECIFIED")).
				AddValue(builder.NewEnumValue("BOOK_VIEW_BASIC")).
				AddValue(builder.NewEnumValue("BOOK_VIEW_FULL"))),
			testutils.Problems{}},
		{"SuffixedView", "DeleteBookRequest", "custom_view", builder.FieldTypeEnum(
			builder.NewEnum("BookView").
				AddValue(builder.NewEnumValue("BOOK_VIEW_UNSPECIFIED")).
				AddValue(builder.NewEnumValue("BOOK_VIEW_BASIC")).
				AddValue(builder.NewEnumValue("BOOK_VIEW_FULL"))),
			testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField("name", builder.FieldTypeString()),
			).AddField(
				builder.NewField(test.fieldName, test.fieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build DeleteBookRequest message.")
			}

			// Run the lint rule, and establish that it returns the correct problems.
			wantProblems := test.problems.SetDescriptor(message.FindFieldByName(test.fieldName))
			gotProblems := unknownFields.Lint(message.GetFile())
			if diff := wantProblems.Diff(gotProblems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
