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

package aip0231

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
			"BatchGetBooksRequest",
			&field{"parent", builder.FieldTypeString()},
			testutils.Problems{},
			nil},
		{
			"MissingField",
			"BatchGetBooksRequest",
			nil,
			testutils.Problems{{Message: "parent"}},
			nil,
		},
		{
			"InvalidType",
			"BatchGetBooksRequest",
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
message BatchGetBooksResponse {
  // Books requested.
  repeated Book books = 1;
}`,
			problems: testutils.Problems{},
		},
		{
			testName: "FieldIsNotRepeated",
			src: `
message BatchGetBooksResponse {
  // Book requested.
  Book book = 1;
}`,
			problems: testutils.Problems{{Message: "The \"Book\" type field on Batch Get Response message should be repeated"}},
		},
		{
			testName: "MissingField",
			src: `
message BatchGetBooksResponse {
  string response = 1;
}`,
			problems: testutils.Problems{{Message: "Message \"BatchGetBooksResponse\" has no \"Book\" type field"}},
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

func TestNamesField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			testName: "Valid-Names",
			src: `
message BatchGetBooksRequest {
repeated string names = 1;
}`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-StandardGetReq",
			src: `
message BatchGetBooksRequest {
repeated GetBookRequest requests = 1;
}

message GetBookRequest {}
`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingNamesField",
			src: `
message BatchGetBooksRequest {
 string parent = 1;
}`,
			problems: testutils.Problems{{Message: `Message "BatchGetBooksRequest" has no "names" field`}},
		},
		{
			testName: "Invalid-KeepingNamesFieldOnly",
			src: `
message BatchGetBooksRequest {
repeated string names = 1;
repeated GetBookRequest requests = 2;
}

message GetBookRequest {}`,
			problems: testutils.Problems{{Message: `Message "BatchGetBooksRequest" should delete "requests" field, only keep the "names" field`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-NamesFieldIsNotRepeated",
			src: `
message BatchGetBooksRequest {
string names = 1;
}`,
			problems: testutils.Problems{{Message: `The "names"" field should be repeated`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("names")
			},
		},
		{
			testName: "Invalid-NamesFieldWrongType",
			src: `
message BatchGetBooksRequest {
repeated int32 names = 1;
}`,
			problems: testutils.Problems{{Message: `"names" field on Batch Get Request should be a "string" type`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("names")
			},
		},
		{
			testName: "Invalid-GetReqFieldIsNotRepeated",
			src: `
message BatchGetBooksRequest {
GetBookRequest requests = 1;
}

message GetBookRequest {}`,
			problems: testutils.Problems{{Message: `The "requests" field should be repeated`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-NamesFieldWrongType",
			src: `
message BatchGetBooksRequest {
	repeated string requests = 1;
}`,
			problems: testutils.Problems{{Message: `The "requests" field on Batch Get Request should be a "GetBookRequest" type`}},
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

			problems := namesField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
