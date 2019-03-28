package visitors

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/reflect/protodesc"
	"github.com/golang/protobuf/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/types/descriptor"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test.protoset --proto_path=testdata testdata/test.proto
func TestScopedMessageVisitor(t *testing.T) {
	f1 := readProtoFile("test.protoset").GetFile()[0]
	fd1, err := protodesc.NewFile(f1, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	tests := []struct {
		name                              protoreflect.Name
		fieldCount, enumCount, oneofCount int
	}{
		{
			name:       "Outer",
			fieldCount: 10,
			enumCount:  2,
			oneofCount: 2,
		},
		{
			name:       "Middle",
			fieldCount: 5,
			enumCount:  1,
			oneofCount: 1,
		},
		{
			name:       "Inner",
			fieldCount: 2,
		},
	}

	for _, test := range tests {
		fieldCount := 0
		enumCount := 0
		oneofCount := 0
		v := &ScopedMessageVisitor{
			Predicate: func(m protoreflect.MessageDescriptor) bool { return m.Name() == test.name },
			Funcs: MessageVisitingFuncs{
				FieldVisit: func(protoreflect.FieldDescriptor) { fieldCount++ },
				EnumVisit:  func(protoreflect.EnumDescriptor) { enumCount++ },
				OneofVisit: func(protoreflect.OneofDescriptor) { oneofCount++ },
			},
		}
		WalkMessage(fd1, v)
		if fieldCount != test.fieldCount {
			t.Errorf("Visiting message %s: got %d fields, but wanted %d", test.name, fieldCount, test.fieldCount)
		}
		if enumCount != test.enumCount {
			t.Errorf("Visiting message %s: got %d enums, but wanted %d", test.name, enumCount, test.enumCount)
		}
		if oneofCount != test.oneofCount {
			t.Errorf("Visiting message %s: got %d oneofs, but wanted %d", test.name, oneofCount, test.oneofCount)
		}
	}
}

func TestScopedMessageVisitor_NoPredicate_VisitAll(t *testing.T) {
	f1 := readProtoFile("test.protoset").GetFile()[0]
	fd1, err := protodesc.NewFile(f1, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	got := 0
	wanted := 10
	v := &ScopedMessageVisitor{
		Funcs: MessageVisitingFuncs{
			FieldVisit: func(protoreflect.FieldDescriptor) { got++ },
		},
	}
	WalkMessage(fd1, v)
	if got != wanted {
		t.Errorf("Visiting message with no predicate: got %d fields, but wanted %d", got, wanted)
	}
}

func TestSimpleMessageVisitor(t *testing.T) {
	f1 := readProtoFile("test.protoset").GetFile()[0]
	fd1, err := protodesc.NewFile(f1, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error: %v", err)
	}

	got := 0
	wanted := 10
	v := &SimpleMessageVisitor{
		Funcs: MessageVisitingFuncs{
			FieldVisit: func(protoreflect.FieldDescriptor) { got++ },
		},
	}
	WalkMessage(fd1, v)
	if got != wanted {
		t.Errorf("Visiting message with no predicate: got %d fields, but wanted %d", got, wanted)
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
