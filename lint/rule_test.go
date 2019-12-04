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
			rule := &FileRule{
				Name: RuleName("test"),
				OnlyIf: func(fd *desc.FileDescriptor) bool {
					return fd.GetName() == "test.proto"
				},
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
			rule := &MessageRule{
				Name: RuleName("test"),
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
			rule := &MessageRule{
				Name: RuleName("test"),
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
			rule := &FieldRule{
				Name: RuleName("test"),
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
				Name: RuleName("test"),
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
				Name: RuleName("test"),
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
				Name: RuleName("test"),
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

func TestEnumValueRule(t *testing.T) {
	// Create a file descriptor with a top-level enum with values.
	fd, err := builder.NewFile("test.proto").AddEnum(
		builder.NewEnum("Format").AddValue(builder.NewEnumValue("YAML")).AddValue(builder.NewEnumValue("JSON")),
	).Build()
	if err != nil {
		t.Fatalf("Catastrophic failure, could not build proto. BOOMZ!")
	}

	for _, test := range makeLintRuleTests(fd.GetEnumTypes()[0].GetValues()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum value rule.
			rule := &EnumValueRule{
				Name: RuleName("test"),
				OnlyIf: func(e *desc.EnumValueDescriptor) bool {
					return e.GetName() == "JSON"
				},
				LintEnumValue: func(e *desc.EnumValueDescriptor) []Problem {
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
				Name: RuleName("test"),
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

func TestDescriptorRule(t *testing.T) {
	// Create a file with one of everything in it.
	book := builder.NewMessage("Book").AddNestedEnum(
		builder.NewEnum("Format").AddValue(
			builder.NewEnumValue("FORMAT_UNSPECIFIED"),
		).AddValue(builder.NewEnumValue("PAPERBACK")),
	).AddField(builder.NewField("name", builder.FieldTypeString())).AddNestedMessage(
		builder.NewMessage("Author"),
	)
	fd, err := builder.NewFile("library.proto").AddMessage(book).AddService(
		builder.NewService("Library").AddMethod(
			builder.NewMethod(
				"ConjureBook",
				builder.RpcTypeMessage(book, false),
				builder.RpcTypeMessage(book, false),
			),
		),
	).AddEnum(builder.NewEnum("State")).Build()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Create a rule that lets us verify that each descriptor was visited.
	visited := make(map[string]desc.Descriptor)
	rule := &DescriptorRule{
		Name: RuleName("test"),
		OnlyIf: func(d desc.Descriptor) bool {
			return d.GetName() != "FORMAT_UNSPECIFIED"
		},
		LintDescriptor: func(d desc.Descriptor) []Problem {
			visited[d.GetName()] = d
			return nil
		},
	}

	// Run the rule.
	rule.Lint(fd)

	// Verify that each descriptor was visited.
	// We do not care what order they were visited in.
	wantDescriptors := []string{
		"Author", "Book", "ConjureBook", "Format", "PAPERBACK",
		"name", "Library", "State",
	}
	if got, want := rule.GetName(), "test"; string(got) != want {
		t.Errorf("Got name %q, wanted %q", got, want)
	}
	if got, want := len(visited), len(wantDescriptors); got != want {
		t.Errorf("Got %d descriptors, wanted %d.", got, want)
	}
	for _, name := range wantDescriptors {
		if _, ok := visited[name]; !ok {
			t.Errorf("Missing descriptor %q.", name)
		}
	}
}

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
			if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], aliases), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from ruleIsEnabled, got %v", want, got)
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
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0], nil), false; got != want {
		t.Errorf("Expected the first message to return %v from ruleIsEnabled, got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[1], nil), true; got != want {
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
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0].GetFields()[0], nil), false; got != want {
		t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[1].GetFields()[0], nil), true; got != want {
		t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
	}
}

type lintRuleTest struct {
	testName string
	problems []Problem
}

// runRule runs a rule within a test environment.
func (test *lintRuleTest) runRule(rule ProtoRule, fd *desc.FileDescriptor, t *testing.T) {
	// Establish that the metadata methods work.
	if got, want := string(rule.GetName()), string(RuleName("test")); got != want {
		t.Errorf("Got %q for GetName(), expected %q", got, want)
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
