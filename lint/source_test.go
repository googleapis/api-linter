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
	"github.com/jgeewax/api-linter/visitors"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_source.protoset --proto_path=testdata testdata/test_source.proto

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

	visitors.WalkMessage(fd1, &visitors.SimpleMessageVisitor{
		Funcs: visitors.MessageVisitingFuncs{
			EnumVisit:      func(f protoreflect.EnumDescriptor) { checkLeadingComment(f, descSource, t) },
			EnumValueVisit: func(f protoreflect.EnumValueDescriptor) { checkLeadingComment(f, descSource, t) },
			FieldVisit:     func(f protoreflect.FieldDescriptor) { checkLeadingComment(f, descSource, t) },
			MessageVisit:   func(f protoreflect.MessageDescriptor) { checkLeadingComment(f, descSource, t) },
			OneofVisit:     func(f protoreflect.OneofDescriptor) { checkLeadingComment(f, descSource, t) },
		},
	})

	visitors.WalkService(fd1, &visitors.SimpleServiceVisitor{
		Funcs: visitors.ServiceVisitFuncs{
			MethodVisit:  func(f protoreflect.MethodDescriptor) { checkLeadingComment(f, descSource, t) },
			ServiceVisit: func(f protoreflect.ServiceDescriptor) { checkLeadingComment(f, descSource, t) },
		},
	})

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
