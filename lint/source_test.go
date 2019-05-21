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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_source.protoset --proto_path=testdata testdata/test_source.proto
//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_rule_disable.protoset --proto_path=testdata testdata/test_rule_disable.proto

type testDescriptorVisiting struct {
	visit func(d protoreflect.Descriptor)
}

func (v testDescriptorVisiting) VisitDescriptor(d protoreflect.Descriptor) {
	v.visit(d)
}

func TestDescriptorLocation(t *testing.T) {
	req := readProtoFile(t, "test_source.protoset")
	fileDesc := req.ProtoFile()
	descSource := req.DescriptorSource()

	tests := []struct {
		descriptor protoreflect.Descriptor
		want       Location
	}{
		{
			descriptor: fileDesc.Messages().Get(0), // A top level message.
			want:       Location{Position{8, 1}, Position{60, 2}},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Messages().Get(0), // A nested message.
			want:       Location{Position{10, 3}, Position{37, 4}},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0),
			want:       Location{Position{45, 3}, Position{50, 4}},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0).Values().Get(1),
			want:       Location{Position{49, 5}, Position{49, 13}},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Fields().Get(1),
			want:       Location{Position{42, 3}, Position{42, 42}},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Oneofs().Get(0),
			want:       Location{Position{54, 3}, Position{59, 4}},
		},
		{
			descriptor: fileDesc.Services().Get(0),
			want:       Location{Position{73, 1}, Position{78, 2}},
		},
		{
			descriptor: fileDesc.Services().Get(0).Methods().Get(0),
			want:       Location{Position{75, 3}, Position{75, 45}},
		},
	}

	for _, test := range tests {
		got, err := descSource.DescriptorLocation(test.descriptor)
		errPrefix := fmt.Sprintf("DescriptorLocation for %s", test.descriptor.FullName())
		if err != nil {
			t.Errorf("%s returns error: %v", errPrefix, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s returns %v, but want %v", errPrefix, got, test.want)
		}
	}
}

func TestDescriptorComments(t *testing.T) {
	req := readProtoFile(t, "test_source.protoset")
	fileDesc := req.ProtoFile()
	descSource := req.DescriptorSource()

	tests := []struct {
		descriptor protoreflect.Descriptor
		want       Comments
	}{
		{
			descriptor: fileDesc.Messages().Get(0),
			want: Comments{
				LeadingComments: " Outer\n",
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0),
			want: Comments{
				LeadingComments: " NestedEnum\n",
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Fields().Get(1),
			want: Comments{
				LeadingComments: " outer_middle_field\n",
			},
		},
	}

	for _, test := range tests {
		got, err := descSource.DescriptorComments(test.descriptor)
		errPrefix := fmt.Sprintf("DescriptorComments for %s", test.descriptor.FullName())
		if err != nil {
			t.Errorf("%s returns error: %v", errPrefix, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s returns %s, but want %s", errPrefix, got, test.want)
		}
	}
}

func TestSyntaxLocation(t *testing.T) {
	req := readProtoFile(t, "test_source.protoset")
	descSource := req.DescriptorSource()
	want := Location{Position{3, 1}, Position{3, 19}}
	got, err := descSource.SyntaxLocation()
	if err != nil {
		t.Errorf("SyntaxLocation() error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SyntaxLocation() returns %v, but want %v", got, want)
	}
}

func TestSyntaxComments(t *testing.T) {
	req := readProtoFile(t, "test_source.protoset")
	descSource := req.DescriptorSource()

	want := Comments{
		LeadingDetachedComments: []string{" DO NOT EDIT -- This is for `source_test.go`.\n"},
	}
	got, err := descSource.SyntaxComments()
	if err != nil {
		t.Errorf("SyntaxComments() error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SyntaxComments() returns %s, but want %s", got, want)
	}
}

func TestIsRuleDisabled(t *testing.T) {
	req := readProtoFile(t, "test_rule_disable.protoset")
	fileDesc := req.ProtoFile()
	descSource := req.DescriptorSource()

	tests := []struct {
		desc     protoreflect.Descriptor
		rule     RuleName
		disabled bool
	}{
		{
			desc:     fileDesc,
			rule:     "rule_disabled_in_file",
			disabled: true,
		},
		{
			desc:     fileDesc.Messages().Get(0),
			rule:     "rule_disabled_in_file",
			disabled: true,
		},
		{
			desc:     fileDesc.Messages().Get(0).Fields().Get(0),
			rule:     "rule_disabled_in_file",
			disabled: true,
		},
		{
			desc:     fileDesc,
			rule:     "rule_not_disabled",
			disabled: false,
		},
		{
			desc:     fileDesc.Messages().Get(0),
			rule:     "rule_not_disabled",
			disabled: false,
		},
		{
			desc:     fileDesc.Messages().Get(0).Fields().Get(0),
			rule:     "rule_not_disabled",
			disabled: false,
		},
		{
			desc:     fileDesc.Messages().Get(0),
			rule:     "rule_disabled_in_message_leading_comment",
			disabled: true,
		},
		{
			desc:     fileDesc.Messages().Get(0).Fields().Get(0),
			rule:     "rule_disabled_in_message_leading_comment",
			disabled: true, // this field is contained by the message, and therefore, the disabling applies here too.
		},
		{
			desc:     fileDesc.Messages().Get(0).Fields().Get(0),
			rule:     "rule_disabled_in_field_trailing_comment",
			disabled: true,
		},
		{
			desc:     fileDesc.Messages().Get(0).Fields().Get(1), // another field
			rule:     "rule_disabled_in_field_trailing_comment",
			disabled: false,
		},
	}

	for _, test := range tests {
		disabled := descSource.isRuleDisabled(test.rule, test.desc)
		if disabled != test.disabled {
			t.Errorf("isRuleDisabled(%s, %s): got %v, but wanted %v", test.rule, test.desc.FullName(), disabled, test.disabled)
		}
	}
}

func readProtoFile(t *testing.T, fileName string) Request {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		t.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	proto := protoset.GetFile()[0]
	req, err := NewProtoRequest(proto)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
