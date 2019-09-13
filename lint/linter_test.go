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
	defaultConfigs := Configs{
		{[]string{"**"}, []string{}, map[string]RuleConfig{}},
	}

	ruleProblems := []Problem{{
		Message:    "rule1_problem",
		Category:   "",
		RuleID:     "test::rule1",
		Descriptor: fd,
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
					RuleConfigs:   map[string]RuleConfig{"": {Disabled: true}},
				},
			),
			ruleProblems,
		},
		{
			"NonMatchingRule",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs:   map[string]RuleConfig{"foo::bar": {Disabled: true}},
				},
			),
			ruleProblems,
		},
		{
			"DisabledRule",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Disabled: true},
					},
				},
			),
			[]Problem{},
		},
		{
			"CategoryOverride",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Category: "error"},
					},
				},
			),
			[]Problem{{
				Category:   "error",
				Descriptor: fd,
				Message:    "rule1_problem",
				RuleID:     "test::rule1",
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			rules, err := NewRuleRegistry(&FileRule{
				Name: "test::rule1",
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

	tests := []struct {
		rule Rule
	}{
		{&FileRule{
			Name: "panic",
			LintFile: func(_ *desc.FileDescriptor) []Problem {
				panic("panic")
			},
		}},
		{&FileRule{
			Name: "panic-error",
			LintFile: func(_ *desc.FileDescriptor) []Problem {
				panic(fmt.Errorf("panic"))
			},
		}},
	}

	for _, test := range tests {
		r, err := NewRuleRegistry(test.rule)
		if err != nil {
			t.Fatalf("Failed to create Rules: %q", err)
		}

		// Instantiate a linter with the given rule.
		l := New(r, []Config{{
			IncludedPaths: []string{"**"},
			RuleConfigs:   map[string]RuleConfig{"": {}},
		}})

		_, err = l.LintProtos(fd)
		if err == nil || !strings.Contains(err.Error(), "panic") {
			t.Fatalf("Expected error with panic, got %q", err)
		}
	}
}
