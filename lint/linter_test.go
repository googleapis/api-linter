package lint

import (
	"reflect"
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

func TestLinter_run(t *testing.T) {
	fileName := "protofile"
	req, _ := NewProtoRequest(
		&descriptorpb.FileDescriptorProto{
			Name: &fileName,
		})

	defaultConfigs := Configs{
		{[]string{"**"}, []string{}, map[string]RuleConfig{}},
	}

	ruleProblems := []Problem{{Message: "rule1_problem", Category: Warning, RuleID: "test::rule1"}}

	tests := []struct {
		desc    string
		configs Configs
		resp    Response
	}{
		{"empty config empty response", Configs{}, Response{FilePath: req.ProtoFile().Path()}},
		{
			"config with non-matching file has no effect",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"nofile"},
					RuleConfigs:   map[string]RuleConfig{"": {Disabled: true}},
				},
			),
			Response{Problems: ruleProblems, FilePath: req.ProtoFile().Path()},
		},
		{
			"config with non-matching rule has no effect",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs:   map[string]RuleConfig{"foo::bar": {Disabled: true}},
				},
			),
			Response{Problems: ruleProblems, FilePath: req.ProtoFile().Path()},
		},
		{
			"matching config can disable rule",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Disabled: true},
					},
				},
			),
			Response{FilePath: req.ProtoFile().Path()},
		},
		{
			"matching config can override Category",
			append(
				defaultConfigs,
				Config{
					IncludedPaths: []string{"*"},
					RuleConfigs: map[string]RuleConfig{
						"test::rule1": {Category: Error},
					},
				},
			),
			Response{
				Problems: []Problem{{Message: "rule1_problem", Category: Error, RuleID: "test::rule1"}},
				FilePath: req.ProtoFile().Path(),
			},
		},
	}

	for ind, test := range tests {
		rules, err := NewRules(&mockRule{
			info:     RuleInfo{Name: "test::rule1"},
			lintResp: ruleProblems,
		})
		if err != nil {
			t.Fatal(err)
		}
		l := New(rules, test.configs)

		resp, _ := l.run(req)
		if !reflect.DeepEqual(resp, test.resp) {
			t.Errorf("Test #%d (%s): Linter.run()=%v; want %v", ind+1, test.desc, resp, test.resp)
		}
	}
}
