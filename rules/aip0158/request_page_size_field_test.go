// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0158

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestRequestPaginationPageSize(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		messageFields []field
		isOneof       bool
		problems      testutils.Problems
		problemDesc   func(m protoreflect.MessageDescriptor) protoreflect.Descriptor
	}{
		{
			"Valid",
			"ListFooRequest",
			[]field{{"page_size", builder.FieldTypeInt32()}, {"page_token", builder.FieldTypeString()}},
			false,
			testutils.Problems{},
			nil,
		},
		{
			"MissingField",
			"ListFooRequest",
			[]field{{"page_token", builder.FieldTypeString()}},
			false,
			testutils.Problems{{Message: "page_size"}},
			nil,
		},
		{
			"InvalidType",
			"ListFooRequest",
			[]field{{"page_size", builder.FieldTypeDouble()}},
			false,
			testutils.Problems{{Suggestion: "int32"}},
			func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.FindFieldByName("page_size")
			},
		},
		{
			"IrrelevantMessage",
			"ListFooPageToken",
			[]field{{"page_token", builder.FieldTypeString()}},
			false,
			nil,
			nil,
		},
		{
			"InvalidIsOneof",
			"ListFooRequest",
			[]field{{"page_size", builder.FieldTypeInt32()}},
			/* isOneof */ true,
			testutils.Problems{{Message: "oneof"}},
			func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.FindFieldByName("page_size")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			for _, f := range test.messageFields {
				fb := builder.NewField(f.fieldName, f.fieldType)
				if test.isOneof {
					messageBuilder.AddOneOf(builder.NewOneOf(fmt.Sprintf("%s_oneof", f.fieldName)).AddChoice(fb))
				} else {
					messageBuilder.AddField(fb)
				}
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Determine the descriptor that a failing test will attach to.
			var problemDesc protoreflect.Descriptor = message
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(message)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := requestPaginationPageSize.Lint(message.ParentFile())
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
