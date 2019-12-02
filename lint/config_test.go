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

func TestReadConfigsJSON(t *testing.T) {
	content := `
	[
		{
			"included_paths": ["path_a"],
			"excluded_paths": ["path_b"],
			"disabled_rules": ["rule_a", "rule_b"],
			"enabled_rules": ["rule_c", "rule_d"]
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
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsJSON returns %v, but want %v", configs, expected)
	}
}
