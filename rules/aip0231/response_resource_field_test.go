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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResourceField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m protoreflect.MessageDescriptor) protoreflect.Descriptor
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
			testName: "ValidEs",
			src: `
				message BatchGetMatchesResponse {
					repeated Match matches = 1;
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
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `
				{{.Src}}
				message Book {}
				message Match {}
			`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Src string }{test.src})

			m := file.Messages().Get(0)

			// Determine the descriptor that a failing test will attach to.
			var problemDesc protoreflect.Descriptor = m.Fields().Get(0)
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := resourceField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
