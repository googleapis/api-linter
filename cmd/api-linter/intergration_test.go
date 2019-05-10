package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Each rule must have a testing case stored here.
// Each case must be positive when the rule in test
// is enabled. It must also contain a "disable-me-here"
// comment at the place where you want the rule to be
// disabled.
var testCases = []struct {
	rule, proto string
}{
	{
		rule: "core::proto_version",
		proto: `
		// disable-me-here
		syntax = "proto2";
		`,
	},
	{
		rule: "core::naming_formats::field_names",
		proto: `
				syntax = "proto3";
				message Test {
					// disable-me-here
					string badName = 1;
				}
			`,
	},
}

func TestEveryRuleHasATestCase(t *testing.T) {
	tests := make(map[string]bool)
	for _, test := range testCases {
		tests[test.rule] = true
	}
	rules := rules()
	for _, rl := range rules {
		ruleID := string(rl.Info().Name)
		if _, found := tests[ruleID]; !found {
			t.Errorf("%s does not have a test case", ruleID)
		}
	}
}

func TestRules_Enabled(t *testing.T) {
	config := `
	[
		{
			"included_paths": ["*.proto"],
			"rule_configs": {
				"": {
					"status": "enabled",
					"category": "warning"
				}
			}
		}
	]
	`

	for _, test := range testCases {
		proto := test.proto
		result := runLinter(t, proto, config)
		if !strings.Contains(result, test.rule) {
			t.Errorf("Rule %q should be enabled by the user config: %q", test.rule, config)
		}
	}
}

func TestRules_DisabledByFileComments(t *testing.T) {
	config := `
	[
		{
			"included_paths": ["*.proto"],
			"rule_configs": {
				"": {
					"status": "enabled",
					"category": "warning"
				}
			}
		}
	]
	`

	for _, test := range testCases {
		disableInFile := fmt.Sprintf("// (-- api-linter: %s=disabled --)\n", test.rule)
		proto := disableInFile + "\n" + test.proto
		result := runLinter(t, proto, config)
		if strings.Contains(result, test.rule) {
			t.Errorf("Rule %q should be disabled by file comments", test.rule)
		}
	}
}

func TestRules_DisabledByInlineComments(t *testing.T) {
	config := `
	[
		{
			"included_paths": ["*.proto"],
			"rule_configs": {
				"": {
					"status": "enabled",
					"category": "warning"
				}
			}
		}
	]
	`

	for _, test := range testCases {
		disableInline := fmt.Sprintf("// (-- api-linter: %s=disabled --)\n", test.rule)
		proto := strings.Replace(test.proto, "disable-me-here", disableInline, -1)
		result := runLinter(t, proto, config)
		if strings.Contains(result, test.rule) {
			t.Errorf("Rule %q should be disabled by in-line comments", test.rule)
		}
	}
}

func TestRules_DisabledByConfig(t *testing.T) {
	config := `
	[
		{
			"included_paths": ["*.proto"],
			"rule_configs": {
				"": {
					"disabled": false,
					"category": "warning"
				}
			}
		},
		{
			"included_paths": ["*.proto"],
			"rule_configs": {
				"replace-me-here": {
					"disabled": true
				}
			}
		}
	]
	`

	for _, test := range testCases {
		c := strings.Replace(config, "replace-me-here", test.rule, -1)
		result := runLinter(t, test.proto, c)
		if strings.Contains(result, test.rule) {
			t.Errorf("Rule %q should be disabled by the user config: %q", test.rule, c)
		}
	}
}

func runLinter(t *testing.T, proto, config string) string {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	protoPath := filepath.Join(workdir, "test.proto")
	configPath := filepath.Join(workdir, "test_config.json")
	outPath := filepath.Join(workdir, "test.out")

	if err := writeFile(protoPath, proto); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(configPath, config); err != nil {
		t.Fatal(err)
	}

	args := []string{
		"api-linter-test",
		"checkproto",
		"--cfg=" + configPath,
		"--out=" + outPath,
		"--proto_path=" + workdir,
		protoPath}
	if err := runCLI(rules(), configs(), args); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func writeFile(path, content string) error {
	if path == "" {
		return nil
	}
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), 0644)
}
