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
		t.Fatalf("Could not build file descriptor.")
	}

	// Declare tests.
	tests := []struct {
		testName string
		problems []Problem
	}{
		{"NoProblems", []Problem{}},
		{"OneProblem", []Problem{{
			Message:    "There was a problem.",
			Descriptor: fd,
		}}},
		{"TwoProblems", []Problem{
			{Message: "This was the first problem.", Descriptor: fd},
			{Message: "This was the second problem.", Descriptor: fd},
		}},
	}

	// Iterate over the tests and run them.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			uri := "https://foo.dev/file-test"
			rule := &FileRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				LintFile: func(fd *desc.FileDescriptor) []Problem {
					return test.problems
				},
			}

			// Ensure that the metadata methods seem correct.
			if got, want := string(rule.GetName()), string(NewRuleName("test", test.testName)); got != want {
				t.Errorf("Got %q for GetName(), expected %q", got, want)
			}
			if got, want := rule.GetURI(), uri; got != want {
				t.Errorf("Got %q for GetURI(), expected %q.", got, want)
			}

			// Run the rule's lint function on the file descriptor
			// and assert that we got what we expect.
			if got, want := rule.Lint(fd), test.problems; !reflect.DeepEqual(got, want) {
				t.Errorf("Got %v problems; expected %v.", got, want)
			}
		})
	}
}

func TestMessageRule(t *testing.T) {
	// Create a file descriptor with two messages in it.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Foo"),
	).AddMessage(
		builder.NewMessage("Bar"),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Declare tests.
	tests := []struct {
		testName string
		problems []Problem
	}{
		{"NoProblems", []Problem{}},
		{"OneProblem", []Problem{{
			Message:    "There was a problem.",
			Descriptor: fd.GetMessageTypes()[1],
		}}},
		{"TwoProblems", []Problem{
			{
				Message:    "This was the first problem.",
				Descriptor: fd.GetMessageTypes()[1],
			},
			{
				Message:    "This was the second problem.",
				Descriptor: fd.GetMessageTypes()[1],
			},
		}},
	}

	// Iterate over the tests and run them.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			uri := "https://foo.dev/message-test"
			rule := &MessageRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				LintMessage: func(m *desc.MessageDescriptor) []Problem {
					if m.GetName() == "Bar" {
						return test.problems
					}
					return nil
				},
			}

			// Establish that the metadata methods work.
			if got, want := string(rule.GetName()), string(NewRuleName("test", test.testName)); got != want {
				t.Errorf("Got %q for GetName(), expected %q", got, want)
			}
			if got, want := rule.GetURI(), uri; got != want {
				t.Errorf("Got %q for GetURI(), expected %q.", got, want)
			}

			// Run the message rule's lint function on the file descriptor
			// and assert that we got what we expect.
			if got, want := rule.Lint(fd), test.problems; !reflect.DeepEqual(got, want) {
				t.Errorf("Got %v problems; expected %v.", got, want)
			}
		})
	}
}

func TestFieldRule(t *testing.T) {
	// Create a file descriptor with two messages in it.
	fd, err := builder.NewFile("test.proto").AddMessage(
		builder.NewMessage("Foo").AddField(
			builder.NewField("bar", builder.FieldTypeString()),
		).AddField(
			builder.NewField("baz", builder.FieldTypeInt32()),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Declare tests.
	problemField := fd.GetMessageTypes()[0].GetFields()[1]
	tests := []struct {
		testName string
		problems []Problem
	}{
		{"NoProblems", []Problem{}},
		{"OneProblem", []Problem{{
			Message:    "There was a problem.",
			Descriptor: problemField,
		}}},
		{"TwoProblems", []Problem{
			{
				Message:    "This was the first problem.",
				Descriptor: problemField,
			},
			{
				Message:    "This was the second problem.",
				Descriptor: problemField,
			},
		}},
	}

	// Iterate over the tests and run them.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			uri := "https://foo.dev/field-test"
			rule := &FieldRule{
				Name: NewRuleName("test", test.testName),
				URI:  uri,
				LintField: func(f *desc.FieldDescriptor) []Problem {
					if f.GetName() == "baz" {
						return test.problems
					}
					return nil
				},
			}

			// Establish that the metadata methods work.
			if got, want := string(rule.GetName()), string(NewRuleName("test", test.testName)); got != want {
				t.Errorf("Got %q for GetName(), expected %q", got, want)
			}
			if got, want := rule.GetURI(), uri; got != want {
				t.Errorf("Got %q for GetURI(), expected %q.", got, want)
			}

			// Run the message rule's lint function on the file descriptor
			// and assert that we got what we expect.
			if got, want := rule.Lint(fd), test.problems; !reflect.DeepEqual(got, want) {
				t.Errorf("Got %v problems; expected %v.", got, want)
			}
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
			if got, want := ruleIsEnabled(rule, f.GetMessageTypes()[0]), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from isEnabled, got %v", want, got)
			}
		})
	}
}
