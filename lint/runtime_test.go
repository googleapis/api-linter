package lint

import (
	"reflect"
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

func TestRepository_Run(t *testing.T) {
	fileName := "protofile"
	req, _ := NewProtoRequest(
		&descriptorpb.FileDescriptorProto{
			Name: &fileName,
		})

	defaultConfigs := RuntimeConfigs{
		{[]string{"**"}, []string{}, map[string]RuleConfig{"": {Status: Enabled}}},
	}

	ruleProblems := []Problem{{Message: "rule1_problem", category: Warning}}

	tests := []struct {
		desc    string
		configs RuntimeConfigs
		resp    Response
	}{
		{"empty config empty response", RuntimeConfigs{}, Response{}},
		{
			"config with non-matching file has no effect",
			append(
				defaultConfigs,
				RuntimeConfig{
					IncludedPaths: []string{"nofile"},
					RuleConfigs:   map[string]RuleConfig{"": {Status: Disabled}},
				},
			),
			Response{Problems: ruleProblems},
		},
		{
			"config with non-matching rule has no effect",
			append(
				defaultConfigs,
				RuntimeConfig{
					IncludedPaths: []string{"*"},
					RuleConfigs:   map[string]RuleConfig{"foo::bar": {Status: Disabled}},
				},
			),
			Response{Problems: ruleProblems},
		},
		{
			"matching config can disable rule",
			append(
				defaultConfigs,
				RuntimeConfig{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Status: Disabled},
					},
				},
			),
			Response{},
		},
		{
			"matching config can override Category",
			append(
				defaultConfigs,
				RuntimeConfig{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Category: Error},
					},
				},
			),
			Response{
				Problems: []Problem{{Message: "rule1_problem", category: Error}},
			},
		},
	}

	for ind, test := range tests {
		runtime := NewRuntime(test.configs)
		err := runtime.AddRules(
			&mockRule{
				info: RuleInfo{Name: "test::rule1"},
				lintResp: Response{
					Problems: ruleProblems,
				},
			})

		if err != nil {
			t.Errorf("Runtime.AddRules(...)=%v; want nil", err)
			continue
		}

		resp, _ := runtime.Run(req)
		if !reflect.DeepEqual(resp, test.resp) {
			t.Errorf("Test #%d (%s): Runtime.Run()=%q; want %q", ind+1, test.desc, resp, test.resp)
		}
	}
}
