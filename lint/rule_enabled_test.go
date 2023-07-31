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

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestRuleIsEnabled(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := &FileRule{
		Name: RuleName("a::b::c"),
		LintFile: func(fd *desc.FileDescriptor) []Problem {
			return []Problem{}
		},
	}

	aliases := map[string]string{
		"a::b::c": "d::e::f",
	}

	// Create appropriate test permutations.
	tests := []struct {
		testName       string
		fileComment    string
		messageComment string
		enabled        bool
	}{
		{"Enabled", "", "", true},
		{"FileDisabled", "api-linter: a::b::c=disabled", "", false},
		{"MessageDisabled", "", "api-linter: a::b::c=disabled", false},
		{"NameNotMatch", "", "api-linter: other=disabled", true},
		{"RegexpNotMatch", "", "api-lint: a::b::c=disabled", true},
		{"AliasDisabled", "", "api-linter: d::e::f=disabled", false},
		{"FileComments_PrefixMatched_Disabled", "api-linter: a=disabled", "", false},
		{"FileComments_MiddleMatched_Disabled", "api-linter: b=disabled", "", false},
		{"FileComments_SuffixMatched_Disabled", "api-linter: c=disabled", "", false},
		{"FileComments_MultipleLinesMatched_Disabled", "api-linter: x=disabled\napi-linter: a=disabled", "", false},
		{"MessageComments_PrefixMatched_Disabled", "", "api-linter: a=disabled", false},
		{"MessageComments_MiddleMatched_Disabled", "", "api-linter: b=disabled", false},
		{"MessageComments_SuffixMatched_Disabled", "", "api-linter: c=disabled", false},
		{"MessageComments_MultipleLinesMatched_Disabled", "", "api-linter: x=disabled\napi-linter: a=disabled", false},
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
			if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], nil, aliases, false), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from ruleIsEnabled, got %v", want, got)
			}
			if !test.enabled {
				if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], nil, aliases, true), true; got != want {
					t.Errorf("Expected the test rule with ignoreCommentDisables true to return %v from ruleIsEnabled, got %v", want, got)
				}
			}
		})
	}
}

func TestRuleIsEnabledFirstMessage(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := &FileRule{
		Name: RuleName("test"),
		LintFile: func(fd *desc.FileDescriptor) []Problem {
			return []Problem{}
		},
	}

	// Build a proto and check that ruleIsEnabled does the right thing.
	f, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("FirstMessage").SetComments(builder.Comments{
			LeadingComment: "api-linter: test=disabled",
		}),
	).AddMessage(
		builder.NewMessage("SecondMessage"),
	).Build()
	if err != nil {
		t.Fatalf("Error building test file: %q", err)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], nil, nil, false), false; got != want {
		t.Errorf("Expected the first message to return %v from ruleIsEnabled, got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[1], nil, nil, false), true; got != want {
		t.Errorf("Expected the second message to return %v from ruleIsEnabled, got %v", want, got)
	}
}

func TestRuleIsEnabledParent(t *testing.T) {
	// Create a rule that we can check enabled status on.
	rule := &FieldRule{
		Name: RuleName("test"),
		LintField: func(f *desc.FieldDescriptor) []Problem {
			return nil
		},
	}

	// Build a proto with two messages, one of which disables the rule.
	f, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Foo").SetComments(builder.Comments{
			LeadingComment: "api-linter: test=disabled",
		}).AddField(builder.NewField("foo", builder.FieldTypeBool())),
	).AddMessage(
		builder.NewMessage("Bar").AddField(builder.NewField("bar", builder.FieldTypeBool())),
	).Build()
	if err != nil {
		t.Fatalf("Error building test file: %q", err)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0].GetFields()[0], nil, nil, false), false; got != want {
		t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[1].GetFields()[0], nil, nil, false), true; got != want {
		t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
	}
}

func TestRuleIsEnabledDeprecated(t *testing.T) {
	// Create a rule that we can check enabled status on.
	rule := &FieldRule{
		Name: RuleName("test"),
		LintField: func(f *desc.FieldDescriptor) []Problem {
			return nil
		},
	}

	for _, test := range []struct {
		name            string
		msgDeprecated   bool
		fieldDeprecated bool
		enabled         bool
	}{
		{"Both", true, true, false},
		{"Message", true, false, false},
		{"Field", false, true, false},
		{"Neither", false, false, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Build a proto with a message and field, possibly deprecated.
			f, err := builder.NewFile("test.proto").AddMessage(
				builder.NewMessage("Foo").SetOptions(&dpb.MessageOptions{
					Deprecated: &test.msgDeprecated,
				}).AddField(builder.NewField("bar", builder.FieldTypeBool()).SetOptions(
					&dpb.FieldOptions{Deprecated: &test.fieldDeprecated},
				)),
			).Build()
			if err != nil {
				t.Fatalf("Error building test file: %q", err)
			}
			if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0].GetFields()[0], nil, nil, false), test.enabled; got != want {
				t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
			}
		})
	}
}
