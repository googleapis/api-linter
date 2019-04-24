package lint

import (
	"reflect"
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

func TestRepository_Run_NoFoundConfig(t *testing.T) {
	fileName := "protofile"
	req, _ := NewProtoFileRequest(
		&descriptorpb.FileDescriptorProto{
			Name: &fileName,
		})
	runtime := NewRuntime(RuleConfig{Status: Enabled, Category: Warning})
	err := runtime.AddRules(
		"test",
		&mockRule{
			info: RuleInfo{Name: "rule1"},
			lintResp: Response{
				Problems: []Problem{{Message: "rule1_problem"}},
			},
		})

	if err != nil {
		t.Errorf("Runtime.AddRules(...)=%v; want nil", err)
		return
	}

	tests := []struct {
		configs Configs
		resp    Response
	}{
		{Configs{}, Response{}},
		{Configs{RuntimeConfig{IncludedPaths: []string{"nofile"}}}, Response{}},
		{
			Configs{
				RuntimeConfig{IncludedPaths: []string{"*"}},
			},
			Response{
				Problems: []Problem{{Message: "rule1_problem", category: Warning}},
			},
		},
		{
			Configs{
				RuntimeConfig{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Status: Disabled},
					},
				},
			},
			Response{},
		},
		{
			Configs{
				RuntimeConfig{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Category: Error},
					},
				},
			},
			Response{
				Problems: []Problem{{Message: "rule1_problem", category: Error}},
			},
		},
	}

	for _, test := range tests {
		resp, _ := runtime.Run(req, test.configs)
		if !reflect.DeepEqual(resp, test.resp) {
			t.Errorf("Runtime.Run returns response %q, but want %q with configs `%v`", resp, test.resp, test.configs)
		}
	}
}
