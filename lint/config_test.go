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
	"strings"
	"testing"
)

func TestRuleConfigs_getRuleConfig(t *testing.T) {
	matchConfig := RuleConfig{Disabled: false, Category: "warning"}

	tests := []struct {
		configs Configs
		path    string
		rule    RuleName
		result  RuleConfig
	}{
		{nil, "a", "b", RuleConfig{}},
		{
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					RuleConfigs: map[string]RuleConfig{
						"foo":      {},
						"testrule": matchConfig,
					},
				},
			},
			"b.proto",
			"testrule",
			RuleConfig{},
		},
		{
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					RuleConfigs: map[string]RuleConfig{
						"foo": matchConfig,
					},
				},
			},
			"a.proto",
			"testrule",
			RuleConfig{},
		},
		{
			Configs{
				{
					IncludedPaths: []string{"a.proto"},
					RuleConfigs: map[string]RuleConfig{
						"foo":      {},
						"testrule": matchConfig,
					},
				},
			},
			"a.proto",
			"testrule",
			matchConfig,
		},
		{
			Configs{
				{
					IncludedPaths: []string{"a/**/*.proto"},
					RuleConfigs: map[string]RuleConfig{
						"foo":     {},
						"a::b::c": matchConfig,
					},
				},
			},
			"a/with/long/sub/dir/ect/ory/e.proto",
			"a::b::c",
			matchConfig,
		},
		{
			Configs{
				{
					IncludedPaths: []string{"a/**/*.proto"},
					RuleConfigs: map[string]RuleConfig{
						"foo":       {},
						"a::module": matchConfig,
					},
				},
			},
			"a/with/long/sub/dir/ect/ory/e.proto",
			"a::module::test_rule",
			matchConfig,
		},
	}
	for ind, test := range tests {
		cfg, _ := test.configs.getRuleConfig(test.path, test.rule)
		if cfg != test.result {
			t.Errorf("Test #%d: %+v.getRuleConfig(%q, %q)=%+v; want %+v", ind+1, test.configs, test.path, test.rule, cfg, test.result)
		}
	}
}

func TestReadConfigsJSON(t *testing.T) {
	content := `
	[
		{
			"included_paths": ["a"],
			"excluded_paths": ["b"],
			"rule_configs": {
				"rule_a": {
					"status": "enabled",
					"category": "warning"
				}
			}
		}
	]
	`

	configs, err := ReadConfigsJSON(strings.NewReader(content))
	if err != nil {
		t.Errorf("ReadConfigsJSON returns error: %v", err)
	}

	expected := Configs{
		{
			IncludedPaths: []string{"a"},
			ExcludedPaths: []string{"b"},
			RuleConfigs: map[string]RuleConfig{
				"rule_a": {
					Disabled: false,
					Category: "warning",
				},
			},
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsJSON returns %v, but want %v", configs, expected)
	}
}

func TestRuleConfig_WithOverride(t *testing.T) {
	tests := []struct {
		original RuleConfig
		override RuleConfig
		result   RuleConfig
	}{
		{
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Disabled: false, Category: "warning"},
		},
		{
			RuleConfig{},
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Disabled: false, Category: "warning"},
		},
		{
			RuleConfig{Category: ""},
			RuleConfig{Disabled: true, Category: "warning"},
			RuleConfig{Disabled: true, Category: "warning"},
		},
		{
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Disabled: true, Category: "error"},
			RuleConfig{Disabled: true, Category: "error"},
		},
		{
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Category: ""},
			RuleConfig{Disabled: false, Category: "warning"},
		},
		{
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Disabled: true, Category: ""},
			RuleConfig{Disabled: true, Category: "warning"},
		},
		{
			RuleConfig{Disabled: false, Category: "warning"},
			RuleConfig{Category: "error"},
			RuleConfig{Disabled: false, Category: "error"},
		},
	}

	for _, test := range tests {
		result := test.original.withOverride(test.override)
		if result != test.result {
			t.Errorf("%+v.WithOverride(%+v)=%+v; want %+v", test.original, test.override, result, test.result)
		}
	}
}
