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
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestInputName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		methodName string
		inputName  string
		msg        string
	}{
		{"Valid", "CreateBook", "CreateBookRequest", ""},
		{"Invalid", "CreateBook", "Book",
			"Post RPCs should have a properly named request message \"CreateBookRequest\", but not \"Book\""},
		{"Irrelevant_OutputWrong", "CreateIamPolicy", "CreateIamPolicyRequest", ""},
		{"Irrelevant_NotCreate", "BuildBook", "Book", ""},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create a minimal service with a AIP-133 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage(test.inputName), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the expected problems.
			problems := inputName.Lint(service.GetFile())
			if test.msg == "" && len(problems) > 0 {
				t.Errorf("Got %v, expected no problems.", problems)
			} else if test.msg != "" && !strings.Contains(problems[0].Message, test.msg) {
				t.Errorf("Got %q, expected message containing %q", problems[0].Message, test.msg)
			}
		})
	}
}
