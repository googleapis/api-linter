package lint

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
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
	fileDesc, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

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
	fileDesc, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

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
	_, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}
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
	_, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

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
	fileDesc, proto := readProtoFile("test_rule_disable.protoset")
	descSource, _ := newDescriptorSource(proto)

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

func readProtoFile(fileName string) (protoreflect.FileDescriptor, *descriptorpb.FileDescriptorProto) {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	proto := protoset.GetFile()[0]
	f, err := protodesc.NewFile(proto, nil)
	if err != nil {
		log.Fatalf("protodesc.NewFile() error: %v", err)
	}
	return f, proto
}
