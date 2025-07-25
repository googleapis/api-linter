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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResourceField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		MessageName string
		FieldName   string
		FieldType   string
		descriptor  func(protoreflect.FileDescriptor) protoreflect.Descriptor
		problems    testutils.Problems
	}{
		{
			"Valid",
			"CreateBookRequest",
			"book",
			"Book",
			nil,
			testutils.Problems{},
		},
		{
			"MissingField",
			"CreateBookRequest",
			"",
			"",
			func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Messages().Get(0)
			},
			testutils.Problems{{Message: "Message \"CreateBookRequest\" has no \"Book\" type field"}},
		},
		{
			"WrongName",
			"CreateBookRequest",
			"payload",
			"Book",
			func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Messages().Get(0).Fields().Get(0)
			},
			testutils.Problems{{Suggestion: "book"}},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `
				message {{.MessageName}} {
					{{.FieldType}} {{.FieldName}} = 1;
				}
				message Book {}
			`
			if test.FieldName == "" {
				template = `
					message {{.MessageName}} {}
					message Book {}
				`
			}
			f := testutils.ParseProto3Tmpl(t, template, test)

			var d protoreflect.Descriptor
			if test.descriptor != nil {
				d = test.descriptor(f)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := resourceField.Lint(f)
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}