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

package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		methodName     string
		reqMessageName string
		problems       testutils.ProblemStubs
	}{
		{"Valid", "GetBook", "GetBookRequest", testutils.ProblemStubs{}},
		{"Invalid", "GetBook", "Book", testutils.ProblemStubs{testutils.ProblemStub{
			Message: "GetBookRequest",
		}}},
		{"GetIamPolicy", "GetIamPolicy", "GetIamPolicyRequest", testutils.ProblemStubs{}},
		{"Irrelevant", "AcquireBook", "Book", testutils.ProblemStubs{}},
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
			method := service.GetMethods()[0]
			test.problems.SetDescriptor(method).Verify(requestMessageName.LintMethod(method), t)
		})
	}
}

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName        string
		methodName      string
		respMessageName string
		problems        testutils.ProblemStubs
	}{
		{"Valid", "GetBook", "Book", testutils.ProblemStubs{}},
		{"Invalid", "GetBook", "GetBookResponse", testutils.ProblemStubs{testutils.ProblemStub{Message: "Book"}}},
		{"Irrelevant", "AcquireBook", "AcquireBookResponse", testutils.ProblemStubs{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-131 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("GetBookRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage(test.respMessageName), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			method := service.GetMethods()[0]
			test.problems.SetDescriptor(method).Verify(responseMessageName.LintMethod(method), t)
		})
	}
}
