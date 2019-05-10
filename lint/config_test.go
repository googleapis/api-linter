package lint

import (
	"reflect"
	"strings"
	"testing"
)

func TestRuleConfigs_getRuleConfig(t *testing.T) {
	matchConfig := RuleConfig{Enabled, Warning}

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
					Status:   "enabled",
					Category: "warning",
				},
			},
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsJSON returns %q, but want %q", configs, expected)
	}
}

func TestRuleConfig_WithOverride(t *testing.T) {
	tests := []struct {
		original RuleConfig
		override RuleConfig
		result   RuleConfig
	}{
		{
			RuleConfig{Enabled, Warning},
			RuleConfig{Enabled, Warning},
			RuleConfig{Enabled, Warning},
		},
		{
			RuleConfig{},
			RuleConfig{Enabled, Warning},
			RuleConfig{Enabled, Warning},
		},
		{
			RuleConfig{Enabled, ""},
			RuleConfig{Disabled, Warning},
			RuleConfig{Disabled, Warning},
		},
		{
			RuleConfig{"", Warning},
			RuleConfig{Disabled, Error},
			RuleConfig{Disabled, Error},
		},
		{
			RuleConfig{Enabled, Warning},
			RuleConfig{"", ""},
			RuleConfig{Enabled, Warning},
		},
		{
			RuleConfig{Enabled, Warning},
			RuleConfig{Disabled, ""},
			RuleConfig{Disabled, Warning},
		},
		{
			RuleConfig{Enabled, Warning},
			RuleConfig{"", Error},
			RuleConfig{Enabled, Error},
		},
	}

	for _, test := range tests {
		result := test.original.withOverride(test.override)
		if result != test.result {
			t.Errorf("%+v.WithOverride(%+v)=%+v; want %+v", test.original, test.override, result, test.result)
		}
	}
}
