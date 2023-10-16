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
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRuleConfigs_IsRuleEnabled(t *testing.T) {
	enabled := true
	disabled := false

	tests := []struct {
		name    string
		configs Configs
		path    string
		rule    string
		want    bool
	}{
		{"EmptyConfig", nil, "a", "b", enabled},
		{
			"NoConfigMatched_Enabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					DisabledRules: []string{"testrule"},
				},
			},
			"b.proto",
			"testrule",
			enabled,
		},
		{
			"PathMatched_DisabledRulesNotMatch_Enabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					DisabledRules: []string{"somerule"},
				},
			},
			"a.proto",
			"testrule",
			enabled,
		},
		{
			"PathExactMatched_DisabledRulesMatched_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					DisabledRules: []string{"somerule", "testrule"},
				},
			},
			"a.proto",
			"testrule",
			disabled,
		},
		{
			"PathExactMatched_DisabledRulesMatchedAll_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					DisabledRules: []string{"all"},
				},
			},
			"a.proto",
			"testrule",
			disabled,
		},
		{
			"PathDoubleStartMatched_DisabledRulesMatched_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a/**/*.proto"},
					DisabledRules: []string{"somerule", "testrule"},
				},
			},
			"a/with/long/sub/dir/etc/ory/e.proto",
			"testrule",
			disabled,
		},
		{
			"PathMatched_DisabledRulesPrefixMatched_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a/b/c.proto"},
					DisabledRules: []string{"parent"},
				},
			},
			"a/b/c.proto",
			"parent::test_rule",
			disabled,
		},
		{
			"PathMatched_DisabledRulesSuffixMatched_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a/b/c.proto"},
					DisabledRules: []string{"child"},
				},
			},
			"a/b/c.proto",
			"parent::child",
			disabled,
		},
		{
			"PathMatched_DisabledRulesMiddleMatched_Disabled",
			Configs{
				{
					IncludedPaths: []string{"a/b/c.proto"},
					DisabledRules: []string{"middle"},
				},
			},
			"a/b/c.proto",
			"parent::middle::child",
			disabled,
		},
		{
			"EmptyIncludePath_ConfigMatched_DisabledRulesMatched_Disabled",
			Configs{
				{
					DisabledRules: []string{"testrule"},
				},
			},
			"a.proto",
			"testrule",
			disabled,
		},
		{
			"ExcludedPathMatch_ConfigNotMatched_DisabledRulesMatched_Enabled",
			Configs{
				{
					ExcludedPaths: []string{"a.proto"},
					DisabledRules: []string{"testrule"},
				},
			},
			"a.proto",
			"testrule",
			enabled,
		},
		{
			"TwoConfigs_Override_Enabled",
			Configs{
				{
					DisabledRules: []string{"testrule"},
				},
				{
					EnabledRules: []string{"testrule::a"},
				},
			},
			"a.proto",
			"testrule::a",
			enabled,
		},
		{
			"TwoConfigs_Override_Disabled",
			Configs{
				{
					EnabledRules: []string{"testrule"},
				},
				{
					DisabledRules: []string{"testrule::a"},
				},
			},
			"a.proto",
			"testrule::a",
			disabled,
		},
		{
			"TwoConfigs_DoubleEnable_Enabled",
			Configs{
				{
					EnabledRules: []string{"testrule"},
				},
				{
					EnabledRules: []string{"testrule::a"},
				},
			},
			"a.proto",
			"testrule::a",
			enabled,
		},
		{
			"TwoConfigs_DoubleDisabled_Disabled",
			Configs{
				{
					DisabledRules: []string{"testrule"},
				},
				{
					DisabledRules: []string{"testrule::a"},
				},
			},
			"a.proto",
			"testrule::a",
			disabled,
		},
		{
			"NoConfigMatched_DefaultDisabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					DisabledRules: []string{"testrule"},
				},
			},
			"b.proto",
			"cloud::25164::generic-fields",
			disabled,
		},
		{
			"ConfigMatched_DefaultDisabled_Enabled",
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					EnabledRules:  []string{"cloud"},
				},
			},
			"a.proto",
			"cloud::25164::generic-fields",
			enabled,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.configs.IsRuleEnabled(test.rule, test.path)

			if got != test.want {
				t.Errorf("IsRuleEnabled: got %t, but want %t", got, test.want)
			}
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("test error")
}

func TestReadConfigsJSONReaderError(t *testing.T) {
	if _, err := ReadConfigsJSON(errReader(0)); err == nil {
		t.Error("ReadConfigsJSON expects an error")
	}
}

func TestReadConfigsJSONFormatError(t *testing.T) {
	invalidJSON := `
	[
		{
			"included_paths": ["path_a"],
			"excluded_paths": ["path_b"],
			"disabled_rules": ["rule_a", "rule_b"],
			"enabled_rules": ["rule_c", "rule_d"]
		}
	`

	if _, err := ReadConfigsJSON(strings.NewReader(invalidJSON)); err == nil {
		t.Error("ReadConfigsJSON expects an error")
	}
}

func TestReadConfigsJSON(t *testing.T) {
	content := `
	[
		{
			"included_paths": ["path_a"],
			"excluded_paths": ["path_b"],
			"disabled_rules": ["rule_a", "rule_b"],
			"enabled_rules": ["rule_c", "rule_d"],
			"import_paths": ["import_a", "import_b"]
		}
	]
	`

	configs, err := ReadConfigsJSON(strings.NewReader(content))
	if err != nil {
		t.Errorf("ReadConfigsJSON returns error: %v", err)
	}

	expected := Configs{
		{
			IncludedPaths: []string{"path_a"},
			ExcludedPaths: []string{"path_b"},
			DisabledRules: []string{"rule_a", "rule_b"},
			EnabledRules:  []string{"rule_c", "rule_d"},
			ImportPaths:   []string{"import_a", "import_b"},
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsJSON returns %v, but want %v", configs, expected)
	}
}

func TestReadConfigsYAMLReaderError(t *testing.T) {
	if _, err := ReadConfigsYAML(errReader(0)); err == nil {
		t.Error("ReadConfigsYAML expects an error")
	}
}

func TestReadConfigsYAMLFormatError(t *testing.T) {
	invalidYAML := `
	[
		{
			"included_paths": ["path_a"],
			"excluded_paths": ["path_b"],
			"disabled_rules": ["rule_a", "rule_b"],
			"enabled_rules": ["rule_c", "rule_d"],
			"import_paths": ["import_a", "import_b"]
		}
	`

	if _, err := ReadConfigsYAML(strings.NewReader(invalidYAML)); err == nil {
		t.Error("ReadConfigsYAML expects an error")
	}
}

func TestReadConfigsYAML(t *testing.T) {
	content := `
---
- included_paths:
    - 'path_a'
  excluded_paths:
    - 'path_b'
  disabled_rules:
    - 'rule_a'
    - 'rule_b'
  enabled_rules:
    - 'rule_c'
    - 'rule_d'
  import_paths:
    - 'import_a'
    - 'import_b'
`

	configs, err := ReadConfigsYAML(strings.NewReader(content))
	if err != nil {
		t.Errorf("ReadConfigsYAML returns error: %v", err)
	}

	expected := Configs{
		{
			IncludedPaths: []string{"path_a"},
			ExcludedPaths: []string{"path_b"},
			DisabledRules: []string{"rule_a", "rule_b"},
			EnabledRules:  []string{"rule_c", "rule_d"},
			ImportPaths:   []string{"import_a", "import_b"},
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsYAML returns %v, but want %v", configs, expected)
	}
}

func TestReadConfigsFromFile(t *testing.T) {
	expectedConfigs := Configs{
		{
			IncludedPaths: []string{"path_a"},
			ExcludedPaths: []string{"path_b"},
			DisabledRules: []string{"rule_a", "rule_b"},
			EnabledRules:  []string{"rule_c", "rule_d"},
			ImportPaths:   []string{"import_a", "import_b"},
		},
	}

	jsonConfigsText := `
	[
		{
			"included_paths": ["path_a"],
			"excluded_paths": ["path_b"],
			"disabled_rules": ["rule_a", "rule_b"],
			"enabled_rules": ["rule_c", "rule_d"],
			"import_paths": ["import_a", "import_b"]
		}
	]
	`
	jsonConfigsFile := createTempFile(t, "test.json", jsonConfigsText)
	defer os.Remove(jsonConfigsFile)

	yamlConfigsText := `
---
- included_paths:
    - 'path_a'
  excluded_paths:
    - 'path_b'
  disabled_rules:
    - 'rule_a'
    - 'rule_b'
  enabled_rules:
    - 'rule_c'
    - 'rule_d'
  import_paths:
    - 'import_a'
    - 'import_b'
`
	yamlConfigsFile := createTempFile(t, "test.yaml", yamlConfigsText)
	defer os.Remove(yamlConfigsFile)

	tests := []struct {
		name     string
		filePath string
		configs  Configs
		hasErr   bool
	}{
		{
			"JSON file",
			jsonConfigsFile,
			expectedConfigs,
			false,
		},
		{
			"YAML file",
			yamlConfigsFile,
			expectedConfigs,
			false,
		},
		{
			"Invalid file extension",
			"test.abc",
			nil,
			true,
		},
		{
			"File not existed",
			"not-existed-file.json",
			nil,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configs, err := ReadConfigsFromFile(test.filePath)
			if (err != nil) != test.hasErr {
				t.Errorf("ReadConfigsFromFile got error %v, but want %v", err, test.hasErr)
			}
			if err != nil {
				if !reflect.DeepEqual(configs, test.configs) {
					t.Errorf("ReadConfigsFromFile got configs %v, but want %v", configs, test.configs)
				}
			}
		})
	}
}

func createTempFile(t *testing.T, name, content string) string {
	dir, err := os.MkdirTemp("", "config_tests")
	if err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dir, name)
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return filePath
}
