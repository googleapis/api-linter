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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResponsePaginationNextPageToken(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		MessageName   string
		Fields        string
		problems      testutils.Problems
		problemDesc   func(m protoreflect.MessageDescriptor) protoreflect.Descriptor
	}{
		{
			"Valid",
			"ListBooksResponse",
			"repeated string books = 1; string next_page_token = 2;",
			testutils.Problems{},
			nil,
		},
		{
			"MissingField",
			"ListBooksResponse",
			"repeated string books = 1;",
			testutils.Problems{{Message: "next_page_token"}},
			nil,
		},
		{
			"InvalidType",
			"ListFooResponse",
			"string name = 1; double next_page_token = 2;",
			testutils.Problems{{Suggestion: "string"}},
			func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("next_page_token")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.Fields}}
				}
			`, test)
			// Determine the descriptor that a failing test will attach to.
			var problemDesc protoreflect.Descriptor = f.Messages().Get(0)
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(f.Messages().Get(0))
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := responsePaginationNextPageToken.Lint(f)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}