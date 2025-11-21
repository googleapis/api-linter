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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

func TestRuleIsEnabled(t *testing.T) {
	// Create a no-op rule, which we can check enabled status on.
	rule := &FileRule{
		Name: RuleName("a::b::c"),
		LintFile: func(fd protoreflect.FileDescriptor) []Problem {
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
			f, err := protodesc.NewFile(&dpb.FileDescriptorProto{
				Name:    proto.String("test.proto"),
				Package: proto.String("test"),
				MessageType: []*dpb.DescriptorProto{
					{
						Name: proto.String("MyMessage"),
					},
				},
				SourceCodeInfo: &dpb.SourceCodeInfo{
					Location: []*dpb.SourceCodeInfo_Location{
						{
							Path:            []int32{2}, // package
							Span:            []int32{1, 1, 1, 1},
							LeadingComments: proto.String(test.fileComment),
						},
						{
							Path:            []int32{4, 0}, // message_type 0
							Span:            []int32{1, 1, 1, 1},
							LeadingComments: proto.String(test.messageComment),
						},
					},
				},
			}, nil)
			if err != nil {
				t.Fatalf("Error building test message: %v", err)
			}
			if got, want := ruleIsEnabled(rule, f.Messages().Get(0), nil, aliases, false), test.enabled; got != want {
				t.Errorf("Expected the test rule to return %v from ruleIsEnabled, got %v", want, got)
			}
			if !test.enabled {
				if got, want := ruleIsEnabled(rule, f.Messages().Get(0), nil, aliases, true), true; got != want {
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
		LintFile: func(fd protoreflect.FileDescriptor) []Problem {
			return []Problem{}
		},
	}

	// Build a proto and check that ruleIsEnabled does the right thing.
	f, err := protodesc.NewFile(&dpb.FileDescriptorProto{
		Name: proto.String("test.proto"),
		MessageType: []*dpb.DescriptorProto{
			{
				Name: proto.String("FirstMessage"),
			},
			{
				Name: proto.String("SecondMessage"),
			},
		},
		SourceCodeInfo: &dpb.SourceCodeInfo{
			Location: []*dpb.SourceCodeInfo_Location{
				{
					Path:            []int32{4, 0}, // message_type 0
					Span:            []int32{1, 1, 1, 1},
					LeadingComments: proto.String("api-linter: test=disabled"),
				},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("Error building test file: %q", err)
	}
	if got, want := ruleIsEnabled(rule, f.Messages().Get(0), nil, nil, false), false; got != want {
		t.Errorf("Expected the first message to return %v from ruleIsEnabled, got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.Messages().Get(1), nil, nil, false), true; got != want {
		t.Errorf("Expected the second message to return %v from ruleIsEnabled, got %v", want, got)
	}
}

func TestRuleIsEnabledParent(t *testing.T) {
	// Create a rule that we can check enabled status on.
	rule := &FieldRule{
		Name: RuleName("test"),
		LintField: func(f protoreflect.FieldDescriptor) []Problem {
			return nil
		},
	}

	// Build a proto with two messages, one of which disables the rule.
	f, err := protodesc.NewFile(&dpb.FileDescriptorProto{
		Name: proto.String("test.proto"),
		MessageType: []*dpb.DescriptorProto{
			{
				Name: proto.String("Foo"),
				Field: []*dpb.FieldDescriptorProto{
					{
						Name:   proto.String("foo"),
						Number: proto.Int32(1),
						Type:   dpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
					},
				},
			},
			{
				Name: proto.String("Bar"),
				Field: []*dpb.FieldDescriptorProto{
					{
						Name:   proto.String("bar"),
						Number: proto.Int32(1),
						Type:   dpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
					},
				},
			},
		},
		SourceCodeInfo: &dpb.SourceCodeInfo{
			Location: []*dpb.SourceCodeInfo_Location{
				{
					Path:            []int32{4, 0}, // message_type 0
					Span:            []int32{1, 1, 1, 1},
					LeadingComments: proto.String("api-linter: test=disabled"),
				},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("Error building test file: %q", err)
	}
	if got, want := ruleIsEnabled(rule, f.Messages().Get(0).Fields().Get(0), nil, nil, false), false; got != want {
		t.Errorf("Expected the foo field to return %v from ruleIsEnabled; got %v", want, got)
	}
	if got, want := ruleIsEnabled(rule, f.Messages().Get(1).Fields().Get(0), nil, nil, false), true; got != want {
		t.Errorf("Expected the bar field to return %v from ruleIsEnabled; got %v", want, got)
	}
}

func TestRuleIsEnabledDeprecated(t *testing.T) {
	// Create a generalRule that we can check enabled status on.
	generalRule := &FieldRule{
		Name: RuleName("test"),
		LintField: func(f protoreflect.FieldDescriptor) []Problem {
			return nil
		},
	}
	deprecationRule := &FieldRule{
		Name: RuleName("core::0192::deprecated-comment"),
		LintField: func(f protoreflect.FieldDescriptor) []Problem {
			return nil
		},
	}

	for _, test := range []struct {
		name            string
		rule            ProtoRule
		msgDeprecated   bool
		fieldDeprecated bool
		enabled         bool
	}{
		{"Both", generalRule, true, true, false},
		{"Message", generalRule, true, false, false},
		{"Field", generalRule, false, true, false},
		{"Neither", generalRule, false, false, true},
		{"DeprecationRule", deprecationRule, false, true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Build a proto with a message and field, possibly deprecated.
			f, err := protodesc.NewFile(&dpb.FileDescriptorProto{
				Name: proto.String("test.proto"),
				MessageType: []*dpb.DescriptorProto{
					{
						Name: proto.String("Foo"),
						Options: &dpb.MessageOptions{
							Deprecated: proto.Bool(test.msgDeprecated),
						},
						Field: []*dpb.FieldDescriptorProto{
							{
								Name:   proto.String("bar"),
								Number: proto.Int32(1),
								Type:   dpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
								Options: &dpb.FieldOptions{
									Deprecated: proto.Bool(test.fieldDeprecated),
								},
							},
						},
					},
				},
			}, nil)
			if err != nil {
				t.Fatalf("Error building test file: %q", err)
			}
			if got, want := ruleIsEnabled(test.rule, f.Messages().Get(0).Fields().Get(0), nil, nil, false), test.enabled; got != want {
				t.Errorf("Expected the bar field to return %v from ruleIsEnabled; got %v", want, got)
			}
		})
	}
}
