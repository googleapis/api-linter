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
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestLinter_run(t *testing.T) {
	fd, err := builder.NewFile("protofile.proto").Build()
	if err != nil {
		t.Fatalf("Failed to build a file descriptor.")
	}
	defaultConfigs := Configs{}

	testRuleName := NewRuleName(111, "test-rule")
	ruleProblems := []Problem{{
		Message:    "rule1_problem",
		Descriptor: fd,
		category:   "",
		RuleID:     testRuleName,
	}}

	tests := []struct {
		testName string
		configs  Configs
		problems []Problem
	}{
		{"Empty", Configs{}, []Problem{}},
		{
			"NonMatchingFile",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"nofile"},
				},
			),
			ruleProblems,
		},
		{
			"NonMatchingRule",
			append(
				defaultConfigs,
				Config{
					DisabledRules: []string{"foo::bar"},
				},
			),
			ruleProblems,
		},
		{
			"DisabledRule",
			append(
				defaultConfigs,
				Config{
					DisabledRules: []string{string(testRuleName)},
				},
			),
			[]Problem{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			rules := NewRuleRegistry()
			err := rules.Register(111, &FileRule{
				Name: NewRuleName(111, "test-rule"),
				LintFile: func(f *desc.FileDescriptor) []Problem {
					return test.problems
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			l := New(rules, test.configs)

			// Actually run the linter.
			resp, _ := l.lintFileDescriptor(fd)

			// Assert that we got the problems we expected.
			if !reflect.DeepEqual(resp.Problems, test.problems) {
				t.Errorf("Expected %v, got %v.", resp.Problems, test.problems)
			}
		})
	}
}

func TestLinter_LintProtos_RulePanics(t *testing.T) {
	fd, err := builder.NewFile("test.proto").Build()
	if err != nil {
		t.Fatalf("Failed to build the file descriptor.")
	}

	testAIP := 111

	tests := []struct {
		testName string
		rule     ProtoRule
	}{
		{
			testName: "Panic",
			rule: &FileRule{
				Name: NewRuleName(testAIP, "panic"),
				LintFile: func(_ *desc.FileDescriptor) []Problem {
					panic("panic")
				},
			},
		},
		{
			testName: "PanicError",
			rule: &FileRule{
				Name: NewRuleName(testAIP, "panic-error"),
				LintFile: func(_ *desc.FileDescriptor) []Problem {
					panic(fmt.Errorf("panic"))
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			rules := NewRuleRegistry()
			err := rules.Register(testAIP, test.rule)
			if err != nil {
				t.Fatalf("Failed to create Rules: %q", err)
			}

			// Instantiate a linter with the given rule.
			l := New(rules, nil)

			_, err = l.LintProtos(fd)
			if err == nil || !strings.Contains(err.Error(), "panic") {
				t.Fatalf("Expected error with panic, got %q", err)
			}
		})
	}
}
