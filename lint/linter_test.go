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
	"github.com/googleapis/api-linter/rules/testutil"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestLinter_run(t *testing.T) {
	fileName := "protofile"
	req, _ := NewProtoRequest(
		&descriptorpb.FileDescriptorProto{
			Name: &fileName,
		}, nil)

	defaultConfigs := Configs{
		{[]string{"**"}, []string{}, map[string]RuleConfig{}},
	}

	ruleProblems := []Problem{{Message: "rule1_problem", Category: "", RuleID: "test::rule1"}}

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
						"test::rule1": {Category: "error"},
					},
				},
			),
			Response{
				Problems: []Problem{{Message: "rule1_problem", Category: "error", RuleID: "test::rule1"}},
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

func TestNewProtoRequest(t *testing.T) {
	depFileDescProto := testutil.MustCreateFileDescriptorProtoFromTemplate(
		"foo.proto",
		`syntax = "proto3";

message Foo {
	string foo_str = 1;
}`,
		nil,
		nil,
	)

	depFileDesc, err := protodesc.NewFile(depFileDescProto, nil)

	if err != nil {
		t.Fatal("Failed to create FileDescriptor: ", err)
	}

	registry := protoregistry.NewFiles(depFileDesc)

	f := testutil.MustCreateFileDescriptorProtoFromTemplate("testfile.proto", `syntax = "proto3";

import "foo.proto";

message Bar {
  Foo foo = 1;
}
`, nil, []*descriptorpb.FileDescriptorProto{depFileDescProto})

	r, err := NewProtoRequest(f, registry)

	if err != nil {
		t.Fatal("NewProtoRequest() returned non nil error: ", err)
	}

	if r.ProtoFile() == nil {
		t.Fatal("r.ProtoFile() = nil; want non-nil protoreflect.FileDescriptor")
	}

	if r.ProtoFile().Messages().Len() != 1 {
		t.Fatalf("r.ProtoFile().Messages().Len()=%d; want 1", r.ProtoFile().Messages().Len())
	}

	barMsg := r.ProtoFile().Messages().Get(0)

	if barMsg.Fields().Len() != 1 {
		t.Fatalf("barMsg.Fields().Len()=%d; want 1", barMsg.Fields().Len())
	}

	fooField := barMsg.Fields().Get(0)

	if fooField.Kind() != protoreflect.MessageKind {
		t.Fatalf("fooField.Kind()=%d; want %d", fooField.Kind(), protoreflect.MessageKind)
	}

	fooMsg := fooField.Message()

	if fooMsg.Fields().Len() != 1 {
		t.Fatalf("fooMsg.Fields().Len()=%d; want 1", fooMsg.Fields().Len())
	}

	fooStrField := fooMsg.Fields().Get(0)

	if fooStrField.Name() != "foo_str" {
		t.Fatalf("fooStrField.Name()=%q; want %q", fooStrField.Name(), "foo_str")
	}
}
