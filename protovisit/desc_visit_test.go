package protovisit

import (
	"testing"

	"github.com/jgeewax/api-linter/protovisit/mocks"
	"github.com/stretchr/testify/mock"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/desc_test.protoset --proto_path=testdata testdata/desc_test.proto

func TestWalkDescriptor(t *testing.T) {
	visitor := new(mocks.DescriptorVisitor)
	visiting := new(mocks.DescriptorVisiting)
	f := readProtoFile("desc_test.protoset")

	// 14 = 4 messages + 3 fields + 2 enums + 2 enum values + 1 oneof + 1 service + 1 method
	numDescriptors := 14
	visitor.On("PreVisit", mock.Anything).Return(nil).Times(numDescriptors)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(numDescriptors)
	visiting.On("VisitDescriptor", mock.Anything).Times(numDescriptors)

	if err := WalkDescriptor(f, visitor, visiting); err != nil {
		t.Errorf("WalkDescriptor: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkDescriptor_ErrSkip(t *testing.T) {
	visitor := new(mocks.DescriptorVisitor)
	visiting := new(mocks.DescriptorVisiting)
	f := readProtoFile("desc_test.protoset")

	// visit the top-level descriptors only, skip and stop immediately.
	// 5 = 3 messages + 1 service + 1 enum.
	numDescriptors := 5
	visitor.On("PreVisit", mock.Anything).Return(ErrSkip).Times(numDescriptors)

	if err := WalkDescriptor(f, visitor, visiting); err != nil {
		t.Errorf("WalkDescriptor: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkDescriptor_ErrSkipVisiting(t *testing.T) {
	visitor := new(mocks.DescriptorVisitor)
	visiting := new(mocks.DescriptorVisiting)
	f := readProtoFile("desc_test.protoset")

	// 14 = 4 messages + 3 fields + 2 enums + 2 enum values + 1 oneof + 1 service + 1 method
	numDescriptors := 14
	visitor.On("PreVisit", mock.Anything).Return(ErrSkipVisiting).Times(numDescriptors)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(numDescriptors)

	if err := WalkDescriptor(f, visitor, visiting); err != nil {
		t.Errorf("WalkDescriptor: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}

func TestWalkDescriptor_ErrSkipNested(t *testing.T) {
	visitor := new(mocks.DescriptorVisitor)
	visiting := new(mocks.DescriptorVisiting)
	f := readProtoFile("desc_test.protoset")

	// visit the top-level descriptors only, skip and stop immediately.
	// 5 = 3 messages + 1 service + 1 enum.
	numDescriptors := 5
	visitor.On("PreVisit", mock.Anything).Return(ErrSkipNested).Times(numDescriptors)
	visitor.On("PostVisit", mock.Anything).Return(nil).Times(numDescriptors)
	visiting.On("VisitDescriptor", mock.Anything).Times(numDescriptors)

	if err := WalkDescriptor(f, visitor, visiting); err != nil {
		t.Errorf("WalkDescriptor: %v", err)
	}
	visitor.AssertExpectations(t)
	visiting.AssertExpectations(t)
}
