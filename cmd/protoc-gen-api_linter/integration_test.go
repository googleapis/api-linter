package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type testFile struct {
	path, content string
}

type testCase struct {
	description           string
	ruleID                string
	protoFile, configFile testFile
	negative              bool
}

var testCases = []testCase{
	{
		description: "positive case",
		ruleID:      "core::proto_version",
		protoFile: testFile{
			path:    "proto_version_test.proto",
			content: `syntax = "proto2";`,
		},
	},
	{
		description: "negative case",
		ruleID:      "core::proto_version",
		protoFile: testFile{
			path:    "proto_version_negative_test.proto",
			content: `syntax = "proto3";`,
		},
		negative: true,
	},
	{
		description: "disabled in config",
		ruleID:      "core::proto_version",
		protoFile: testFile{
			path:    "proto_version_test.proto",
			content: `syntax = "proto2";`,
		},
		configFile: testFile{
			path: "proto_version_config.json",
			content: `[{
				"included_paths": ["**/*.proto"],
				"rule_configs": {
					"core::proto_version": {
						"status": "disabled"
					}
				}
			}]`,
		},
		negative: true,
	},
	{
		description: "positive case",
		ruleID:      "core::naming_formats::field_names",
		protoFile: testFile{
			path: "field_names_test.proto",
			content: `
				syntax = "proto3";
				message Test {
					string badName = 1;
				}
			`,
		},
	},
	{
		description: "negative case",
		ruleID:      "core::naming_formats::field_names",
		protoFile: testFile{
			path: "field_names_test.proto",
			content: `
				syntax = "proto3";
				message Test {
					string good_name = 1;
				}
			`,
		},
		negative: true,
	},
	{
		description: "disabled by in file comment",
		ruleID:      "core::naming_formats::field_names",
		protoFile: testFile{
			path: "field_names_test.proto",
			content: `
				syntax = "proto3";
				message Test {
					// (--api-linter core::naming_formats::field_names=disabled --)
					string good_name = 1;
				}
			`,
		},
		negative: true,
	},
	{
		description: "disabled by config",
		ruleID:      "core::naming_formats::field_names",
		protoFile: testFile{
			path: "field_names_test.proto",
			content: `
				syntax = "proto3";
				message Test {
					string badName = 1;
				}
			`,
		},
		configFile: testFile{
			path: "field_names_config.json",
			content: `[{
				"included_paths": ["**/*.proto"],
				"rule_configs": {
					"core::naming_formats::field_names": {
						"status": "disabled"
					}
				}
			}]`,
		},
		negative: true,
	},
}

func TestEveryRuleHasPositiveTestCase(t *testing.T) {
	rules := getRules()
	positiveTests := make(map[string]testCase)
	for _, test := range testCases {
		if !test.negative {
			positiveTests[test.ruleID] = test
		}
	}
	for _, rl := range rules {
		ruleID := string(rl.Info().Name)
		if _, found := positiveTests[ruleID]; !found {
			t.Errorf("%s does not have a positive test case", ruleID)
		}
	}
}

func TestRules(t *testing.T) {
	for _, test := range testCases {
		result := runLinter(t, test.protoFile, test.configFile)
		if got, want := strings.Contains(result, test.ruleID), !test.negative; got != want {
			t.Errorf("%s: rule %q in linting result %q: %v, but want %v -- proto %q, config %q",
				test.description, test.ruleID, result, got, want, test.protoFile.content, test.configFile.content)
		}
	}
}

func TestRulesDisabledInFile(t *testing.T) {
	for _, test := range testCases {
		disableComment := fmt.Sprintf("(-- api-linter: %s=disabled --)", test.ruleID)
		test.protoFile.content = "//" + disableComment + "\n" + test.protoFile.content
		result := runLinter(t, test.protoFile, test.configFile)
		if strings.Contains(result, test.ruleID) {
			t.Errorf("rule %q is disabled on %q, but got result %q", test.ruleID, test.protoFile.content, result)
		}
	}
}

func runLinter(t *testing.T, protoFile, configFile testFile) string {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	if err := writeFile(workdir, protoFile); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(workdir, configFile); err != nil {
		t.Fatal(err)
	}

	params := "out_file=test.out"
	if configFile.path != "" {
		params += ",cfg_file=" + filepath.Join(workdir, configFile.path)
	}
	runProtoC(t, workdir, params, filepath.Join(workdir, protoFile.path))

	f, err := os.Open(filepath.Join(workdir, "test.out"))
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

func writeFile(rootDir string, file testFile) error {
	if file.path == "" {
		return nil
	}
	path := filepath.Join(rootDir, file.path)
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(file.content), 0644)
}

func runProtoC(t *testing.T, workdir, params string, args ...string) {
	cmd := exec.Command("protoc", "--plugin=protoc-gen-api_linter="+os.Args[0])
	cmd.Args = append(cmd.Args, "--api_linter_out="+params+":"+workdir, "-I"+workdir)
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(os.Environ(), "RUN_AS_PROTOC_PLUGIN=1")
	out, err := cmd.CombinedOutput()

	if err != nil && len(out) > 0 {
		t.Log("RUNNING: ", strings.Join(cmd.Args, " "))
		t.Log(string(out))
	}
	if err != nil {
		t.Fatalf("protoc: %v", err)
	}
}

func init() {
	if os.Getenv("RUN_AS_PROTOC_PLUGIN") != "" {
		run()
		os.Exit(0)
	}
}
