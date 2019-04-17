package rules

import (
	"fmt"
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/lint"
)

func TestCheckFieldNamesUseLowerSnakeCase(t *testing.T) {
	tests := []struct {
		fieldName  string
		numProblem int
		suggestion string
	}{
		{"good_field_name", 0, ""},
		{"badFieldName", 1, "bad_field_name"},
		{"badField_Name", 1, "bad_field_name"},
		{"bad_Field_Name", 1, "bad_field_name"},
	}
	for _, test := range tests {
		proto := &descriptorpb.FileDescriptorProto{
			MessageType: []*descriptorpb.DescriptorProto{
				&descriptorpb.DescriptorProto{
					Field: []*descriptorpb.FieldDescriptorProto{
						&descriptorpb.FieldDescriptorProto{
							Name: &test.fieldName,
						},
					},
				},
			},
		}
		req, _ := lint.NewProtoFileRequest(proto)
		field := req.ProtoFile().Messages().Get(0).Fields().Get(0)

		problems, err := checkFieldNamesUseLowerSnakeCase(field, req.DescriptorSource())

		testErrPrefix := fmt.Sprintf("checkFieldNamesUseLowerSnakeCase for the field with name '%s'", test.fieldName)
		if err != nil {
			t.Errorf("%s returns unexpected error: %v", testErrPrefix, err)
		}
		if test.numProblem != len(problems) {
			t.Errorf("%s returns %d problems, but want %d", testErrPrefix, len(problems), test.numProblem)
		}
		if len(problems) > 0 && problems[0].Suggestion != test.suggestion {
			t.Errorf("%s returns suggestion '%s', but want '%s'", testErrPrefix, problems[0].Suggestion, test.suggestion)
		}
	}
}
