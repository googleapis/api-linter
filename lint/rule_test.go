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

package lint

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestIsEnabled(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := Rule{Name: NewRuleName("test")}

	// Create appropriate test permutations.
	tests := []struct {
		testName       string
		fileComment    string
		messageComment string
		enabled        bool
	}{
		{"Enabled", "", "", true},
		{"FileDisabled", "api-linter: test=disabled", "", false},
		{"MessageDisabled", "", "api-linter: test=disabled", false},
	}

	// Run the specific tests individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, err := builder.NewFile("test.proto").SetSyntaxComments(builder.Comments{
				LeadingComment: test.fileComment,
			}).AddMessage(
				builder.NewMessage("MyMessage").SetComments(builder.Comments{
					LeadingComment: test.messageComment,
				}),
			).Build()
			if err != nil {
				t.Fatalf("Error building test message")
			}
			if got, want := rule.isEnabled(f.GetMessageTypes()[0]), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from isEnabled, got %v", want, got)
			}
		})
	}
}
