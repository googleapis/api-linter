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

package aip0233

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

type field struct {
	fieldName string
	fieldType *builder.FieldType
}

func TestParentField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName     string
		messageName  string
		messageField *field
		problems     testutils.Problems
		problemDesc  func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			"Valid",
			"BatchCreateBooksRequest",
			&field{"parent", builder.FieldTypeString()},
			testutils.Problems{},
			nil},
		{
			"MissingField",
			"BatchCreateBooksRequest",
			nil,
			testutils.Problems{{Message: "parent"}},
			nil,
		},
		{
			"InvalidType",
			"BatchCreateBooksRequest",
			&field{"parent", builder.FieldTypeDouble()},
			testutils.Problems{{Message: "string"}},
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("parent")
			},
		},
		{
			"Irrelevant",
			"EnumerateBooksRequest",
			&field{"id", builder.FieldTypeString()},
			testutils.Problems{},
			nil,
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			if test.messageField != nil {
				messageBuilder.AddField(
					builder.NewField(test.messageField.fieldName, test.messageField.fieldType),
				)
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = message
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(message)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := parentField.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestRequestsField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			testName: "Valid",
			src: `
message BatchCreateBooksRequest {
	repeated CreateBookRequest requests = 1;
}

message CreateBookRequest {}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingRequestsField",
			src: `
message BatchCreateBooksRequest {
	string parent = 1;
}`,
			problems: testutils.Problems{{Message: `Message "BatchCreateBooksRequest" has no "requests" field`}},
		},
		{
			testName: "Invalid-RequestsFieldIsNotRepeated",
			src: `
message BatchCreateBooksRequest {
	CreateBookRequest requests = 1;
}
message CreateBookRequest {}`,
			problems: testutils.Problems{{Message: `The "requests" field should be repeated`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsFieldWrongType",
			src: `
message BatchCreateBooksRequest {
	repeated int32 requests = 1;
}`,
			problems: testutils.Problems{{Message: `The "requests" field on Batch Create Request should be a "CreateBookRequest" type`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsNotRepeatedWrongType",
			src: `
message BatchCreateBooksRequest {
	int32 requests = 1;
}`,
			problems: testutils.Problems{
				{Message: `The "requests" field should be repeated`},
				{Message: `The "requests" field on Batch Create Request should be a "CreateBookRequest" type`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)

			m := file.GetMessageTypes()[0]

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = m
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := requestsField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestResourceField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			testName: "Valid",
			src: `
message BatchCreateBooksResponse {
 // Books requested.
 repeated Book books = 1;
}`,
			problems: testutils.Problems{},
		},
		{
			testName: "FieldIsNotRepeated",
			src: `
message BatchCreateBooksResponse {
 // Book requested.
 Book book = 1;
}`,
			problems: testutils.Problems{{Message: "The \"Book\" type field on Batch Create Response message should be repeated"}},
		},
		{
			testName: "MissingField",
			src: `
message BatchCreateBooksResponse {
 string response = 1;
}`,
			problems: testutils.Problems{{Message: "Message \"BatchCreateBooksResponse\" has no \"Book\" type field"}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `
{{.Src}}
message Book {
}`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Src string }{test.src})

			m := file.GetMessageTypes()[0]

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = m.GetFields()[0]
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := resourceField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
