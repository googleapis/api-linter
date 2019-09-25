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

package aip0158

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

type field struct {
	fieldName string
	fieldType *builder.FieldType
}

func TestRequestPaginationPageSize(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		messageFields []field
		problemCount  int
		errPrefix     string
	}{
		{"Valid", "ListFooRequest", []field{{"page_size", builder.FieldTypeInt32()}, {"page_token", builder.FieldTypeString()}}, 0, "False positive"},
		{"MissingField", "ListFooRequest", []field{{"page_token", builder.FieldTypeString()}}, 1, "False positive"},
		{"InvalidType", "ListFooRequest", []field{{"page_size", builder.FieldTypeDouble()}}, 1, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			for _, f := range test.messageFields {
				messageBuilder.AddField(
					builder.NewField(f.fieldName, f.fieldType),
				)
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := requestPaginationPageSize.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, requestPaginationPageSize.Name, problems)
			}
		})
	}
}

func TestRequestPaginationPageToken(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		messageFields []field
		problemCount  int
		errPrefix     string
	}{
		{"Valid", "ListFooRequest", []field{{"page_token", builder.FieldTypeString()}}, 0, "False positive"},
		{"MissingField", "ListFooRequest", []field{{"name", builder.FieldTypeString()}}, 1, "False positive"},
		{"InvalidType", "ListFooRequest", []field{{"page_token", builder.FieldTypeDouble()}}, 1, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			for _, f := range test.messageFields {
				messageBuilder.AddField(
					builder.NewField(f.fieldName, f.fieldType),
				)
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := requestPaginationPageToken.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, requestPaginationPageToken.Name, problems)
			}
		})
	}
}

func TestResponsePaginationNextPageToken(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		messageFields []field
		problemCount  int
		errPrefix     string
	}{
		{"Valid", "ListFooResponse", []field{{"name", builder.FieldTypeString()}, {"next_page_token", builder.FieldTypeString()}}, 0, "False positive"},
		{"MissingField", "ListFooResponse", []field{{"name", builder.FieldTypeString()}}, 1, "False positive"},
		{"InvalidType", "ListFooResponse", []field{{"name", builder.FieldTypeString()}, {"next_page_token", builder.FieldTypeDouble()}}, 1, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			for _, f := range test.messageFields {
				messageBuilder.AddField(
					builder.NewField(f.fieldName, f.fieldType),
				)
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := responsePaginationNextPageToken.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, responsePaginationNextPageToken.Name, problems)
			}
		})
	}
}
