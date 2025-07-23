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

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestFileRule(t *testing.T) {
	// Create a file descriptor with nothing in it.
	fd := buildFile(t, `syntax = "proto3";`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd) {
		t.Run(test.testName, func(t *testing.T) {
			rule := &FileRule{
				Name: RuleName("test"),
				OnlyIf: func(fd protoreflect.FileDescriptor) bool {
					return fd.Path() == "test.proto"
				},
				LintFile: func(fd protoreflect.FileDescriptor) []Problem {
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
	fd := buildFile(t, `syntax = "proto3"; message Book {} message Author {}`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.Messages().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			rule := &MessageRule{
				Name: RuleName("test"),
				OnlyIf: func(m protoreflect.MessageDescriptor) bool {
					return m.Name() == "Author"
				},
				LintMessage: func(m protoreflect.MessageDescriptor) []Problem {
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
	fd := buildFile(t, `syntax = "proto3"; message Book { message Author {} }`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.Messages().Get(0).Messages().Get(0)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			rule := &MessageRule{
				Name: RuleName("test"),
				OnlyIf: func(m protoreflect.MessageDescriptor) bool {
					return m.Name() == "Author"
				},
				LintMessage: func(m protoreflect.MessageDescriptor) []Problem {
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
	fd := buildFile(t, `syntax = "proto3"; message Book { string title = 1; int32 edition_count = 2; }`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.Messages().Get(0).Fields().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the field rule.
			rule := &FieldRule{
				Name: RuleName("test"),
				OnlyIf: func(f protoreflect.FieldDescriptor) bool {
					return f.Name() == "edition_count"
				},
				LintField: func(f protoreflect.FieldDescriptor) []Problem {
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
	fd := buildFile(t, `syntax = "proto3"; service Library {}`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.Services().Get(0)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the service rule.
			rule := &ServiceRule{
				Name: RuleName("test"),
				LintService: func(s protoreflect.ServiceDescriptor) []Problem {
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
	fd := buildFile(t, `
		syntax = "proto3";
		message Book {}
		message GetBookRequest {}
		message CreateBookRequest {}
		service Library {
			rpc GetBook(GetBookRequest) returns (Book);
			rpc CreateBook(CreateBookRequest) returns (Book);
		}
	`)

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.Services().Get(0).Methods().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the method rule.
			rule := &MethodRule{
				Name: RuleName("test"),
				OnlyIf: func(m protoreflect.MethodDescriptor) bool {
					return m.Name() == "CreateBook"
				},
				LintMethod: func(m protoreflect.MethodDescriptor) []Problem {
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
	fd := buildFile(t, `
		syntax = "proto3";
		enum Format { PDF = 0; }
		enum Edition { PUBLISHER_ONLY = 0; }
	`)

	for _, test := range makeLintRuleTests(fd.Enums().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: RuleName("test"),
				OnlyIf: func(e protoreflect.EnumDescriptor) bool {
					return e.Name() == "Edition"
				},
				LintEnum: func(e protoreflect.EnumDescriptor) []Problem {
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
	fd := buildFile(t, `
		syntax = "proto3";
		enum Format { YAML = 0; JSON = 1; }
	`)

	for _, test := range makeLintRuleTests(fd.Enums().Get(0).Values().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum value rule.
			rule := &EnumValueRule{
				Name: RuleName("test"),
				OnlyIf: func(e protoreflect.EnumValueDescriptor) bool {
					return e.Name() == "JSON"
				},
				LintEnumValue: func(e protoreflect.EnumValueDescriptor) []Problem {
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
	fd := buildFile(t, `
		syntax = "proto3";
		message Book {
			enum Format { PDF = 0; }
			enum Edition { PUBLISHER_ONLY = 0; }
		}
	`)

	for _, test := range makeLintRuleTests(fd.Messages().Get(0).Enums().Get(1)) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: RuleName("test"),
				OnlyIf: func(e protoreflect.EnumDescriptor) bool {
					return e.Name() == "Edition"
				},
				LintEnum: func(e protoreflect.EnumDescriptor) []Problem {
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
	fd := buildFile(t, `
		syntax = "proto3";
		message Book {
			enum Format {
				FORMAT_UNSPECIFIED = 0;
				PAPERBACK = 1;
			}
			message Author {}
			string name = 1;
		}
		service Library {
			rpc ConjureBook(Book) returns (Book);
		}
		enum State {
			AVAILABLE = 0;
		}
	`)

	// Create a rule that lets us verify that each descriptor was visited.
	visited := make(map[string]protoreflect.Descriptor)
	rule := &DescriptorRule{
		Name: RuleName("test"),
		OnlyIf: func(d protoreflect.Descriptor) bool {
			return d.Name() != "FORMAT_UNSPECIFIED"
		},
		LintDescriptor: func(d protoreflect.Descriptor) []Problem {
			visited[string(d.Name())] = d
			return nil
		},
	}

	// Run the rule.
	rule.Lint(fd)

	// Verify that each descriptor was visited.
	// We do not care what order they were visited in.
	wantDescriptors := []string{
		"Author", "Book", "ConjureBook", "Format", "PAPERBACK",
		"name", "Library", "State", "AVAILABLE",
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

type lintRuleTest struct {
	testName string
	problems []Problem
}

// runRule runs a rule within a test environment.
func (test *lintRuleTest) runRule(rule ProtoRule, fd protoreflect.FileDescriptor, t *testing.T) {
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
func makeLintRuleTests(d protoreflect.Descriptor) []lintRuleTest {
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
