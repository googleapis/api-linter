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
	"fmt"
	"reflect"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestFileRule(t *testing.T) {
	// Create a file descriptor with nothing in it.
	fd, err := builder.NewFile("test.proto").Build()
	if err != nil {
		t.Fatalf("Could not build file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd) {
		t.Run(test.testName, func(t *testing.T) {
			uri := fmt.Sprintf("https://aip.dev/%s", t.Name())
			rule := &FileRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				LintFile: func(fd *desc.FileDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestMessageRule(t *testing.T) {
	// Create a file descriptor with two messages in it.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Book"),
	).AddMessage(
		builder.NewMessage("Author"),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			uri := fmt.Sprintf("https://aip.dev/%s", t.Name())
			rule := &MessageRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				OnlyIf: func(m *desc.MessageDescriptor) bool {
					return m.GetName() == "Author"
				},
				LintMessage: func(m *desc.MessageDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

// Establish that nested messages are tested.
func TestMessageRuleNested(t *testing.T) {
	// Create a file descriptor with a message and nested message in it.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Book").AddNestedMessage(builder.NewMessage("Author")),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetNestedMessageTypes()[0]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			uri := fmt.Sprintf("https://aip.dev/%s", t.Name())
			rule := &MessageRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				OnlyIf: func(m *desc.MessageDescriptor) bool {
					return m.GetName() == "Author"
				},
				LintMessage: func(m *desc.MessageDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestFieldRule(t *testing.T) {
	// Create a file descriptor with one message and two fields in that message.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Book").AddField(
			builder.NewField("title", builder.FieldTypeString()),
		).AddField(
			builder.NewField("edition_count", builder.FieldTypeInt32()),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetFields()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			uri := fmt.Sprintf("https://aip.dev/%s", t.Name())
			rule := &FieldRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				OnlyIf: func(f *desc.FieldDescriptor) bool {
					return f.GetName() == "edition_count"
				},
				LintField: func(f *desc.FieldDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestServiceRule(t *testing.T) {
	// Create a file descriptor with a service.
	fd, err := builder.NewFile("test.proto").AddService(
		builder.NewService("Library"),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build a file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetServices()[0]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the service rule.
			rule := &ServiceRule{
				Name: NewRuleName("test", test.testName),
				URI:  fmt.Sprintf("https://aip.dev/%s", t.Name()),
				LintService: func(s *desc.ServiceDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestMethodRule(t *testing.T) {
	// Create a file descriptor with a service.
	book := builder.RpcTypeMessage(builder.NewMessage("Book"), false)
	fd, err := builder.NewFile("test.proto").AddService(
		builder.NewService("Library").AddMethod(
			builder.NewMethod(
				"GetBook",
				builder.RpcTypeMessage(builder.NewMessage("GetBookRequest"), false),
				book,
			),
		).AddMethod(
			builder.NewMethod(
				"CreateBook",
				builder.RpcTypeMessage(builder.NewMessage("CreateBookRequest"), false),
				book,
			),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build a file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetServices()[0].GetMethods()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the method rule.
			rule := &MethodRule{
				Name: NewRuleName("test", test.testName),
				URI:  fmt.Sprintf("https://aip.dev/%s", t.Name()),
				OnlyIf: func(m *desc.MethodDescriptor) bool {
					return m.GetName() == "CreateBook"
				},
				LintMethod: func(m *desc.MethodDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestEnumRule(t *testing.T) {
	// Create a file descriptor with top-level enums.
	fd, err := builder.NewFile("test.proto").AddEnum(
		builder.NewEnum("Format"),
	).AddEnum(
		builder.NewEnum("Edition"),
	).Build()
	if err != nil {
		t.Fatalf("Catastrophic failure, could not build proto. BOOMZ!")
	}

	for _, test := range makeLintRuleTests(fd.GetEnumTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: NewRuleName("test", test.testName),
				URI:  fmt.Sprintf("https://aip.dev/%s", t.Name()),
				OnlyIf: func(e *desc.EnumDescriptor) bool {
					return e.GetName() == "Edition"
				},
				LintEnum: func(e *desc.EnumDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestEnumRuleNested(t *testing.T) {
	// Create a file descriptor with top-level enums.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Book").AddNestedEnum(
			builder.NewEnum("Format"),
		).AddNestedEnum(
			builder.NewEnum("Edition"),
		),
	).Build()
	if err != nil {
		t.Fatalf("Catastrophic failure, could not build proto. BOOMZ!")
	}

	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetNestedEnumTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: NewRuleName("test", test.testName),
				URI:  fmt.Sprintf("https://aip.dev/%s", t.Name()),
				OnlyIf: func(e *desc.EnumDescriptor) bool {
					return e.GetName() == "Edition"
				},
				LintEnum: func(e *desc.EnumDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestRuleIsEnabled(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := &FileRule{
		Name: NewRuleName("test"),
		LintFile: func(fd *desc.FileDescriptor) []Problem {
			return []Problem{}
		},
	}

	aliases := map[string]string{
		"test": "alias",
	}

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
		{"NameNotMatch", "", "api-linter: other=disabled", true},
		{"AliasDisabled", "", "api-linter: alias=disabled", false},
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
			if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], aliases), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from ruleIsEnabled, got %v", want, got)
			}
		})
	}
}

func TestRuleIsEnabledFirstMessage(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := &FileRule{
		Name: NewRuleName("test"),
		LintFile: func(fd *desc.FileDescriptor) []Problem {
			return []Problem{}
		},
	}

	// Run the specific tests individually.
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
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], nil), false; got != want {
		t.Errorf("Expected the first message to return %v from ruleIsEnabled, got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[1], nil), true; got != want {
		t.Errorf("Expected the second message to return %v from ruleIsEnabled, got %v", want, got)
	}
}

type lintRuleTest struct {
	testName string
	problems []Problem
}

// runRule runs a rule within a test environment.
func (test *lintRuleTest) runRule(rule protoRule, fd *desc.FileDescriptor, t *testing.T) {
	// Establish that the metadata methods work.
	if got, want := string(rule.GetName()), string(NewRuleName("test", test.testName)); got != want {
		t.Errorf("Got %q for GetName(), expected %q", got, want)
	}
	if got, want := rule.GetURI(), fmt.Sprintf("https://aip.dev/%s", t.Name()); got != want {
		t.Errorf("Got %q for GetURI(), expected %q.", got, want)
	}

	// Run the rule's lint function on the file descriptor
	// and assert that we got what we expect.
	if got, want := rule.Lint(fd), test.problems; !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v problems; expected %v.", got, want)
	}
}

// makeLintRuleTests generates boilerplate tests that are consistent for
// each type of rule.
func makeLintRuleTests(d desc.Descriptor) []lintRuleTest {
	return []lintRuleTest{
		{"NoProblems", []Problem{}},
		{"OneProblem", []Problem{{
			Message:    "There was a problem.",
			Descriptor: d,
		}}},
		{"TwoProblems", []Problem{
			{
				Message:    "This was the first problem.",
				Descriptor: d,
			},
			{
				Message:    "This was the second problem.",
				Descriptor: d,
			},
		}},
	}
}
