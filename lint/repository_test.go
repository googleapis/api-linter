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
	repo := NewRepository()
	repo.AddRule(
		"test",
		RuleConfig{Status: Enabled, Category: Warning},
		&mockRule{
			info: RuleInfo{Name: "rule1"},
			lintResp: Response{
				Problems: []Problem{{Message: "rule1_problem"}},
			},
		})

	tests := []struct {
		configs Configs
		resp    Response
	}{
		{Configs{}, Response{}},
		{Configs{Config{IncludedPaths: []string{"nofile"}}}, Response{}},
		{
			Configs{
				Config{IncludedPaths: []string{"*"}},
			},
			Response{
				Problems: []Problem{{Message: "rule1_problem", category: Warning}},
			},
		},
		{
			Configs{
				Config{
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
				Config{
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
		resp, _ := repo.Run(req, test.configs)
		if !reflect.DeepEqual(resp, test.resp) {
			t.Errorf("Repository.Run returns response %q, but want %q with configs `%v`", resp, test.resp, test.configs)
		}
	}
}
