package rules

import (
	"fmt"
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/lint"
)

func TestCheckProtoSyntaxVersion(t *testing.T) {
	tests := []struct {
		syntax     string
		numProblem int
	}{
		{"proto3", 0},
		{"proto2", 1},
	}
	for _, test := range tests {
		proto := &descriptorpb.FileDescriptorProto{
			Syntax: &test.syntax,
		}
		req, _ := lint.NewProtoFileRequest(proto)

		problems, err := checkProtoSyntaxVersion(req.ProtoFile(), req.DescriptorSource())

		testErrPrefix := fmt.Sprintf("checkProtoSyntaxVersion for syntax '%s'", test.syntax)
		if err != nil {
			t.Errorf("%s returns unexpected error: %v", testErrPrefix, err)
		}
		if test.numProblem != len(problems) {
			t.Errorf("%s returns %d problems, but want %d", testErrPrefix, len(problems), test.numProblem)
		}
	}
}
