package lint

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/protovisit"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_source.protoset --proto_path=testdata testdata/test_source.proto
//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_rule_disable.protoset --proto_path=testdata testdata/test_rule_disable.proto

type testDescriptorVisiting struct {
	visit func(d protoreflect.Descriptor)
}

func (v testDescriptorVisiting) VisitDescriptor(d protoreflect.Descriptor) {
	v.visit(d)
}

func TestSourceDescriptor(t *testing.T) {
	proto := readProtoFile("test_source.protoset").GetFile()[0]
	f, err := protodesc.NewFile(proto, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	s, err := NewDescriptorSource(proto)
	if err != nil {
		t.Errorf("NewDescriptorSource: %v", err)
	}

	protovisit.WalkDescriptor(
		f,
		protovisit.SimpleDescriptorVisitor{},
		testDescriptorVisiting{
			visit: func(d protoreflect.Descriptor) {
				checkLeadingComment(d, s, t)
			},
		},
	)
}

func TestIsRuleDisabled(t *testing.T) {
	proto := readProtoFile("test_rule_disable.protoset").GetFile()[0]
	f, err := protodesc.NewFile(proto, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	s, err := NewDescriptorSource(proto)
	if err != nil {
		t.Errorf("NewDescriptorSource: %v", err)
	}

	tests := []struct {
		rule  RuleID
		count int
	}{
		{
			rule:  RuleID{Set: "core", Name: "rule_all_disabled"},
			count: 2,
		},
		{
			rule:  RuleID{Set: "core", Name: "rule_not_disabled"},
			count: 0,
		},
		{
			rule:  RuleID{Set: "other", Name: "rule_not_disabled"},
			count: 0,
		},
		{
			rule:  RuleID{Set: "core", Name: "rule_leading_disabled"},
			count: 1,
		},
		{
			rule:  RuleID{Set: "core", Name: "rule_trailing_disabled"},
			count: 1,
		},
		{
			rule:  RuleID{Set: "other", Name: "rule_leading_disabled"},
			count: 1,
		},
	}

	for _, test := range tests {
		count := 0
		protovisit.WalkDescriptor(
			f,
			protovisit.SimpleDescriptorVisitor{},
			testDescriptorVisiting{
				visit: func(d protoreflect.Descriptor) {
					if s.IsRuleDisabled(test.rule, d) {
						count++
					}
				},
			},
		)
		if count != test.count {
			t.Errorf("IsRuleDisabled: got %d for %s, but wanted %d", count, test.rule, test.count)
		}
	}
}

func checkLeadingComment(f protoreflect.Descriptor, descSource DescriptorSource, t *testing.T) {
	comments, err := descSource.DescriptorComments(f)
	if err != nil {
		t.Errorf("DescriptorComments for `%s`: %v", f.FullName(), err)
	}
	leadingComment := strings.TrimSpace(comments.LeadingComments)
	if leadingComment != string(f.Name()) {
		t.Errorf("DescriptorComments for `%s`: got '%s', but wanted '%s'", f.FullName(), leadingComment, f.Name())
	}
}

func readProtoFile(fileName string) *descriptorpb.FileDescriptorSet {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	return protoset
}
