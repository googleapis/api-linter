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

package aip0132

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		methodName     string
		reqMessageName string
		problemCount   int
		errPrefix      string
	}{
		{"Valid", "ListBooks", "ListBooksRequest", 0, "False positive"},
		{"Invalid", "ListBooks", "Books", 1, "False negative"},
		{"Irrelevant", "EnumerateBooks", "Books", 0, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-131 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage(test.reqMessageName), false),
				builder.RpcTypeMessage(builder.NewMessage("ListBooksResponse"), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := requestMessageName.LintMethod(service.GetMethods()[0]); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, requestMessageName.Name, problems)
			}
		})
	}
}

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName        string
		methodName      string
		respMessageName string
		problemCount    int
		errPrefix       string
	}{
		{"Valid", "ListBooks", "ListBooksResponse", 0, "False positive"},
		{"Invalid", "ListBooks", "Books", 1, "False negative"},
		{"Irrelevant", "EnumerateBooks", "EnumerateBooksResponse", 0, "False positive"},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-131 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("ListBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage(test.respMessageName), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			if problems := responseMessageName.LintMethod(service.GetMethods()[0]); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, responseMessageName.Name, problems)
			}
		})
	}
}
