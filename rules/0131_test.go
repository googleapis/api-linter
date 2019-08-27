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

package rules

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestCheckGetRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		methodName     string
		reqMessageName string
		problemCount   int
		errPrefix      string
	}{
		{"Valid", "GetBook", "GetBookRequest", 0, "False positive"},
		{"Invalid", "GetBook", "Book", 1, "False negative"},
		{"Irrelevant", "AcquireBook", "Book", 0, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-131 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage(test.reqMessageName), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := checkGetRequestMessageName.LintMethod(service.GetMethods()[0]); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, checkGetRequestMessageName.Name, problems)
			}
		})
	}
}

func TestCheckGetRequestMessageNameField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		nameFieldName string
		nameFieldType *builder.FieldType
		problemCount  int
		errPrefix     string
	}{
		{"Valid", "GetBookRequest", "name", builder.FieldTypeString(), 0, "False positive"},
		{"InvalidName", "GetBookRequest", "id", builder.FieldTypeString(), 1, "False negative"},
		{"InvalidType", "GetBookRequest", "name", builder.FieldTypeBytes(), 1, "False negative"},
		{"Irrelevant", "AcquireBookRequest", "id", builder.FieldTypeString(), 0, "False positive"},
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
			if problems := checkGetRequestMessageNameField.LintMessage(message); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, checkGetRequestMessageNameField.Name, problems)
			}
		})
	}
}
