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

type testMessageVisitor struct {
	descSource DescriptorSource
	t          *testing.T
}

func (v testMessageVisitor) VisitEnum(f protoreflect.EnumDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testMessageVisitor) VisitEnumValue(f protoreflect.EnumValueDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testMessageVisitor) VisitField(f protoreflect.FieldDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testMessageVisitor) VisitOneof(f protoreflect.OneofDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testMessageVisitor) VisitMessage(f protoreflect.MessageDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testMessageVisitor) VisitExtension(f protoreflect.ExtensionDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

type testServiceVisitor struct {
	descSource DescriptorSource
	t          *testing.T
}

func (v testServiceVisitor) VisitService(f protoreflect.ServiceDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func (v testServiceVisitor) VisitMethod(f protoreflect.MethodDescriptor) {
	checkLeadingComment(f, v.descSource, v.t)
}

func TestSourceDescriptor(t *testing.T) {
	f1 := readProtoFile("test_source.protoset").GetFile()[0]
	fd1, err := protodesc.NewFile(f1, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	descSource, err := NewDescriptorSource(f1)
	if err != nil {
		t.Errorf("NewDescriptorSource: %v", err)
	}

	protovisit.WalkMessage(fd1, protovisit.SimpleMessageVisitor{}, testMessageVisitor{descSource, t})

	protovisit.WalkService(fd1, protovisit.SimpleServiceVisitor{}, testServiceVisitor{descSource, t})

	for i := 0; i < fd1.Enums().Len(); i++ {
		e := fd1.Enums().Get(i)
		checkLeadingComment(e, descSource, t)
		for j := 0; j < e.Values().Len(); j++ {
			checkLeadingComment(e.Values().Get(i), descSource, t)
		}
	}
}

func checkLeadingComment(f protoreflect.Descriptor, descSource DescriptorSource, t *testing.T) {
	comments, err := descSource.FindCommentsByDescriptor(f)
	if err != nil {
		t.Errorf("FindCommentsByDescriptor for `%s`: %v", f.FullName(), err)
	}
	leadingComment := strings.TrimSpace(comments.LeadingComments)
	if leadingComment != string(f.Name()) {
		t.Errorf("FindCommentsByDescriptor for `%s`: got '%s', but wanted '%s'", f.FullName(), leadingComment, f.Name())
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
