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
)

func TestLint(t *testing.T) {
	ruleName1 := NewRuleName(111, "a")
	ruleName2 := NewRuleName(111, "b")
	tests := []struct {
		name        string
		descriptors []Descriptor
		rules       []Rule
		responses   []Response
	}{
		{
			"NoRule",
			[]Descriptor{&mockDescriptor{}},
			[]Rule{},
			[]Response{},
		},
		{
			"OneRule",
			[]Descriptor{&mockDescriptor{}},
			[]Rule{&mockRule{ruleName1}},
			[]Response{{
				FilePath: "",
				Problems: []Problem{{RuleID: ruleName1, Message: string(ruleName1)}},
			}},
		},
		{
			"TwoRule",
			[]Descriptor{&mockDescriptor{}},
			[]Rule{&mockRule{ruleName1}, &mockRule{ruleName2}},
			[]Response{
				{
					FilePath: "",
					Problems: []Problem{
						{RuleID: ruleName1, Message: string(ruleName1)},
						{RuleID: ruleName2, Message: string(ruleName2)},
					},
				},
			},
		},
		{
			"OneRule_TwoFiles",
			[]Descriptor{
				&mockDescriptor{sourceInfo: mockSourceInfo{file: mockFileInfo{path: "test-file1"}}},
				&mockDescriptor{sourceInfo: mockSourceInfo{file: mockFileInfo{path: "test-file2"}}},
			},
			[]Rule{&mockRule{ruleName1}},
			[]Response{
				{
					FilePath: "test-file1",
					Problems: []Problem{{RuleID: ruleName1, Message: string(ruleName1)}},
				},
				{
					FilePath: "test-file2",
					Problems: []Problem{{RuleID: ruleName1, Message: string(ruleName1)}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules := NewRuleRegistry()
			err := rules.Register(111, test.rules...)
			if err != nil {
				t.Fatal(err)
			}
			linter := New(rules, Configs{})
			responses := linter.Lint(test.descriptors...)
			if !reflect.DeepEqual(responses, test.responses) {
				t.Errorf("Lint got %v, but want %v", responses, test.responses)
			}
		})
	}
}

func TestLint_Configs(t *testing.T) {
	ruleName := NewRuleName(111, "test")
	rules := NewRuleRegistry()
	err := rules.Register(111, &mockRule{name: ruleName})
	if err != nil {
		t.Fatal(err)
	}
	descriptor := mockDescriptor{
		sourceInfo: mockSourceInfo{
			file: mockFileInfo{
				path: "test-file",
			},
		},
	}

	tests := []struct {
		name        string
		configs     Configs
		wantProblem bool
	}{
		{"EmptyConfigs_Run", Configs{}, true},
		{"PathNotMatched_Run", Configs{{IncludedPaths: []string{"nofile"}}}, true},
		{
			"PathIncluded_RuleDisabled_NotRun",
			Configs{{
				IncludedPaths: []string{"test-file"},
				DisabledRules: []string{string(ruleName)},
			}},
			false,
		},
		{
			"PathIncluded_RuleEnabled_Run",
			Configs{{
				IncludedPaths: []string{"test-file"},
				EnabledRules:  []string{string(ruleName)},
			}},
			true,
		},
		{
			"PathExcluded_RuleDisabled_Run",
			Configs{{
				ExcludedPaths: []string{"test-file"},
				DisabledRules: []string{string(ruleName)},
			}},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			linter := New(rules, test.configs)
			problems := linter.Lint(descriptor)
			gotProblem := len(problems) > 0
			if test.wantProblem != gotProblem {
				t.Errorf("Running lint? Got %v, but want %v", gotProblem, test.wantProblem)
			}
		})
	}
}

func TestLint_CommentDisabling(t *testing.T) {
	ruleName := NewRuleName(111, "test")
	rules := NewRuleRegistry()
	err := rules.Register(111, &mockRule{name: ruleName})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		descriptor  Descriptor
		wantProblem bool
	}{
		{"NoComments", mockDescriptor{}, true},
		{
			"Disabled_LeadingComments",
			mockDescriptor{
				sourceInfo: mockSourceInfo{leadingComments: fmt.Sprintf("api-linter: %s=disabled", ruleName)},
			},
			false,
		},
		{
			"Disabled_FileComments",
			mockDescriptor{
				sourceInfo: mockSourceInfo{file: mockFileInfo{comments: fmt.Sprintf("api-linter: %s=disabled", ruleName)}},
			},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			linter := New(rules, Configs{})
			problems := linter.Lint(test.descriptor)
			gotProblem := len(problems) > 0
			if test.wantProblem != gotProblem {
				t.Errorf("Running lint? Got %v, but want %v", gotProblem, test.wantProblem)
			}
		})
	}
}

type mockRule struct {
	name RuleName
}

func (r mockRule) Name() RuleName {
	return r.name
}

func (r mockRule) Lint(d Descriptor) []Problem {
	return []Problem{Problem{Message: string(r.name)}}
}

type mockFileInfo struct {
	path, comments string
}

func (f mockFileInfo) Path() string {
	return f.path
}

func (f mockFileInfo) Comments() string {
	return f.comments
}

type mockSourceInfo struct {
	leadingComments string
	file            mockFileInfo
}

func (s mockSourceInfo) LeadingComments() string {
	return s.leadingComments
}

func (s mockSourceInfo) File() FileInfo {
	return s.file
}

type mockDescriptor struct {
	sourceInfo mockSourceInfo
}

func (d mockDescriptor) SourceInfo() SourceInfo {
	return d.sourceInfo
}
