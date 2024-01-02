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

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Each case must be positive when the rule in test
// is enabled. It must also contain a "disable-me-here"
// comment at the place where you want the rule to be
// disabled.
var testCases = []struct {
	testName, rule, proto string
}{
	{
		testName: "GetRequestMessage",
		rule:     "core::0131::request-message-name",
		proto: `
		syntax = "proto3";

		service Library {
			// disable-me-here
			rpc GetBook(Book) returns (Book);
		}

		message Book {}
		`,
	},
	{
		testName: "PackageVersion",
		rule:     "core::0215::versioned-packages",
		proto: `
		syntax = "proto3";

		// disable-me-here
		package google.test;

		message Test {}
		`,
	},
	{
		testName: "FieldNames",
		rule:     "core::0140::lower-snake",
		proto: `
				syntax = "proto3";
				import "dummy.proto";
				message Test {
					// disable-me-here
					string badName = 1;
					dummy.Dummy dummy = 2;
				}
			`,
	},
}

func TestRules_EnabledByDefault(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			proto := test.proto
			result := runLinter(t, proto, "")
			if !strings.Contains(result, test.rule) {
				t.Errorf("rule %q should be enabled by default", test.rule)
			}
		})
	}
}

func TestRules_DisabledByFileComments(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			disableInFile := fmt.Sprintf("// (-- api-linter: %s=disabled --)", test.rule)
			proto := disableInFile + "\n" + test.proto
			result := runLinter(t, proto, "")
			if strings.Contains(result, test.rule) {
				t.Errorf("rule %q should be disabled by file comments", test.rule)
			}
		})
	}
}

func TestRules_DisabledByInlineComments(t *testing.T) {
	testConfigurations := []struct {
		suffix       string
		appendArgs   []string
		wantDisabled bool
	}{
		{
			suffix:       "WithIgnoreCommentDisablesTrue",
			appendArgs:   []string{"--ignore-comment-disables=true"},
			wantDisabled: false,
		},
		{
			suffix:       "WithIgnoreCommentDisablesFalse",
			appendArgs:   []string{"--ignore-comment-disables=false"},
			wantDisabled: true,
		},
		{
			suffix:       "WithNoFlags",
			appendArgs:   []string{},
			wantDisabled: true,
		},
	}

	for _, test := range testCases {
		for _, testConfig := range testConfigurations {
			t.Run(test.testName+testConfig.suffix, func(t *testing.T) {
				disableInline := fmt.Sprintf("(-- api-linter: %s=disabled --)", test.rule)
				proto := strings.Replace(test.proto, "disable-me-here", disableInline, -1)
				_, result := runLinterWithFailureStatus(t, proto, "", testConfig.appendArgs)
				isDisabled := !strings.Contains(result, test.rule)
				if isDisabled != testConfig.wantDisabled {
					t.Errorf("want %q disabled by in-line comments to be %v with flags %v, got %v", test.rule, testConfig.wantDisabled, testConfig.appendArgs, isDisabled)
				}
			})
		}
	}
}

func TestRules_DisabledByConfig(t *testing.T) {
	config := `
	[
		{
			"disabled_rules": ["replace-me-here"]
		}
	]
	`

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			c := strings.Replace(config, "replace-me-here", test.rule, -1)
			result := runLinter(t, test.proto, c)
			if strings.Contains(result, test.rule) {
				t.Errorf("rule %q should be disabled by the user config: %q", test.rule, c)
			}
		})
	}
}

func TestBuildErrors(t *testing.T) {
	expected := `internal/testdata/build_errors.proto:8:1: expected ';'
internal/testdata/build_errors.proto:13:1: expected ';'`
	err := runCLI([]string{"internal/testdata/build_errors.proto"})
	if err == nil {
		t.Fatal("expected build error for build_errors.proto")
	}
	actual := err.Error()
	if expected != actual {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestExitStatusForLintFailure(t *testing.T) {
	type testCase struct{ testName, rule, proto string }
	failCase := testCase{
		testName: "GetRequestMessage",
		rule:     "core::0131::request-message-name",
		proto: `
		syntax = "proto3";

		service Library {
			// disable-me-here
			rpc GetBook(Book) returns (Book);
		}

		message Book {}
		`,
	}
	// checks lint failure = true when lint problems found
	t.Run(failCase.testName+"ReturnsFailure", func(t *testing.T) {
		lintFailureStatus, result := runLinterWithFailureStatus(t, failCase.proto, "", []string{"--set-exit-status"})
		if lintFailureStatus == false {
			t.Log(result)
			t.Fatalf("Expected: %v Actual: %v", true, lintFailureStatus)
		}
	})

	// checks lint failure = false when no problems found
	t.Run(failCase.testName+"ReturnsNoFailure", func(t *testing.T) {
		disableAll := `[ { "disabled_rules": [ "all" ] } ]`
		lintFailureStatus, result := runLinterWithFailureStatus(t, failCase.proto, disableAll, []string{"--set-exit-status"})
		if lintFailureStatus {
			t.Log(result)
		}
		if lintFailureStatus == true {
			t.Fatalf("Expected: %v Actual: %v", false, lintFailureStatus)
		}
	})

	// checks lint failure = false when lint problems found but --set-exit-status not set
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			proto := test.proto
			lintFailureStatus, result := runLinterWithFailureStatus(t, proto, "", []string{})
			expected := result == ""
			if lintFailureStatus != expected {
				t.Fatalf("Expected: %v Actual: %v", expected, lintFailureStatus)
			}
		})
	}
}

func runLinter(t *testing.T, protoContent, configContent string) string {
	_, result := runLinterWithFailureStatus(t, protoContent, configContent, []string{})
	return result
}

func runLinterWithFailureStatus(t *testing.T, protoContent, configContent string, appendArgs []string) (bool, string) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Prepare command line flags.
	args := []string{}
	// Add a flag for the linter config file if the provided
	// config content is not empty.
	if configContent != "" {
		configFileName := "test_config.json"
		configFilePath := filepath.Join(tempDir, configFileName)
		if err := writeFile(configFilePath, configContent); err != nil {
			t.Fatal(err)
		}
		args = append(args, fmt.Sprintf("--config=%s", configFilePath))
	}
	// Add a flag for the output path.
	outPath := filepath.Join(tempDir, "test.out")
	args = append(args, fmt.Sprintf("-o=%s", outPath))
	// Add the temp dir to the proto paths.
	args = append(args, fmt.Sprintf("-I=%s", tempDir))
	// Add a flag for the file descriptor set.
	args = append(args, "--descriptor-set-in=internal/testdata/dummy.protoset")
	// Write the proto file.
	protoFileName := "test.proto"
	protoFilePath := filepath.Join(tempDir, protoFileName)
	if err := writeFile(protoFilePath, protoContent); err != nil {
		t.Fatal(err)
	}
	args = append(args, protoFileName)
	args = append(args, appendArgs...)

	lintErr := runCLI(args)
	if lintErr != nil && !errors.Is(lintErr, ExitForLintFailure) {
		t.Fatal(lintErr)
	}

	out, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	return errors.Is(lintErr, ExitForLintFailure), string(out)
}

func writeFile(path, content string) error {
	if path == "" {
		return nil
	}
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
