package protovisit

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/protovisit/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/golang/protobuf/v2/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test.protoset --proto_path=testdata testdata/test.proto
//go:generate mockery -all

func TestWalkEnum(t *testing.T) {
	visitor := new(mocks.EnumVisitor)
	visiting := new(mocks.EnumVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(nil).Twice()
	visitor.On("PostVisit", mock.Anything).Return(nil).Twice()
	visiting.On("VisitEnum", mock.Anything).Twice()
	visiting.On("VisitEnumValue", mock.Anything).Twice()

	if err := WalkEnum(f, visitor, visiting); err != nil {
		t.Errorf("WalkEnum: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkEnum_ErrSkip(t *testing.T) {
	visitor := new(mocks.EnumVisitor)
	visiting := new(mocks.EnumVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkip).Twice()

	if err := WalkEnum(f, visitor, visiting); err != nil {
		t.Errorf("WalkEnum: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkEnum_ErrSkipVisiting(t *testing.T) {
	visitor := new(mocks.EnumVisitor)
	visiting := new(mocks.EnumVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkipVisiting).Twice()
	visitor.On("PostVisit", mock.Anything).Return(nil).Twice()

	if err := WalkEnum(f, visitor, visiting); err != nil {
		t.Errorf("WalkEnum: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkEnum_Err(t *testing.T) {
	visitor := new(mocks.EnumVisitor)
	visiting := new(mocks.EnumVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(errors.New("some error")).Once()

	if err := WalkEnum(f, visitor, visiting); err == nil {
		t.Error("WalkEnum: expecting error, but got none")
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage(t *testing.T) {
	visitor := new(mocks.MessageVisitor)
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(nil).Times(7)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(7)
	visiting.On("VisitMessage", mock.Anything).Times(7)
	visiting.On("VisitField", mock.Anything).Times(10)
	visiting.On("VisitOneof", mock.Anything).Times(2)
	visiting.On("VisitEnum", mock.Anything).Times(2)
	visiting.On("VisitEnumValue", mock.Anything).Times(4)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage_Err(t *testing.T) {
	visitor := new(mocks.MessageVisitor)
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(errors.New("some error")).Once()

	if err := WalkMessage(f, visitor, visiting); err == nil {
		t.Error("WalkMessage: expecting error, but got none")
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage_ErrSkip(t *testing.T) {
	visitor := new(mocks.MessageVisitor)
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkip).Times(5)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage_ErrSkipVisiting(t *testing.T) {
	visitor := new(mocks.MessageVisitor)
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkipVisiting).Times(7)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(7)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage_ErrSkipNested(t *testing.T) {
	visitor := new(mocks.MessageVisitor)
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkipNested).Times(5)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(5)
	visiting.On("VisitMessage", mock.Anything).Times(5)
	visiting.On("VisitField", mock.Anything).Times(5)
	visiting.On("VisitOneof", mock.Anything).Times(1)
	visiting.On("VisitEnum", mock.Anything).Times(1)
	visiting.On("VisitEnumValue", mock.Anything).Times(2)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkMessage_ScopedVisitor(t *testing.T) {
	visitor := &ScopedMessageVisitor{
		Predicate: func(m protoreflect.MessageDescriptor) bool {
			return m.Name() == "Middle"
		},
	}
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visiting.On("VisitMessage", mock.Anything).Times(2)
	visiting.On("VisitField", mock.Anything).Times(5)
	visiting.On("VisitOneof", mock.Anything).Times(1)
	visiting.On("VisitEnum", mock.Anything).Times(1)
	visiting.On("VisitEnumValue", mock.Anything).Times(2)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visiting.AssertExpectations(t)
}

func TestWalkMessage_ScopedVisitor_NoPredicate(t *testing.T) {
	visitor := &ScopedMessageVisitor{}
	visiting := new(mocks.MessageVisiting)
	f := readProtoFile("test.protoset")

	visiting.On("VisitMessage", mock.Anything).Times(7)
	visiting.On("VisitField", mock.Anything).Times(10)
	visiting.On("VisitOneof", mock.Anything).Times(2)
	visiting.On("VisitEnum", mock.Anything).Times(2)
	visiting.On("VisitEnumValue", mock.Anything).Times(4)

	if err := WalkMessage(f, visitor, visiting); err != nil {
		t.Errorf("WalkMessage: %v", err)
	}
	visiting.AssertExpectations(t)
}

func TestWalkService(t *testing.T) {
	visitor := new(mocks.ServiceVisitor)
	visiting := new(mocks.ServiceVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(nil).Times(1)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(1)
	visiting.On("VisitService", mock.Anything).Once()
	visiting.On("VisitMethod", mock.Anything).Twice()

	if err := WalkService(f, visitor, visiting); err != nil {
		t.Errorf("WalkService: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkService_Err(t *testing.T) {
	visitor := new(mocks.ServiceVisitor)
	visiting := new(mocks.ServiceVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(errors.New("some error")).Times(1)

	if err := WalkService(f, visitor, visiting); err == nil {
		t.Error("WalkService: expecting error, but got none")
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkService_ErrSkip(t *testing.T) {
	visitor := new(mocks.ServiceVisitor)
	visiting := new(mocks.ServiceVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkip).Times(1)

	if err := WalkService(f, visitor, visiting); err != nil {
		t.Errorf("WalkService: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkService_ErrSkipVisiting(t *testing.T) {
	visitor := new(mocks.ServiceVisitor)
	visiting := new(mocks.ServiceVisiting)
	f := readProtoFile("test.protoset")

	visitor.On("PreVisit", mock.Anything).Return(ErrSkipVisiting).Times(1)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(1)

	if err := WalkService(f, visitor, visiting); err != nil {
		t.Errorf("WalkService: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func readProtoFile(fileName string) protoreflect.FileDescriptor {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	f, err := protodesc.NewFile(protoset.GetFile()[0], nil)
	if err != nil {
		log.Fatalf("protodesc.NewFile() error: %v", err)
	}
	return f
}
