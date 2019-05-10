package corp

import (
	"testing"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/lint"
)

func TestProtoFilesMustIncludeVersion(t *testing.T) {
	rule := protoFilesMustIncludeVersion()

	tests := []struct {
		path        string
		numProblems int
	}{
		{"a/b/c/v1/foo.proto", 0},
		{"google/corp/a/b/c/d/e/f/alpha/foo.proto", 0},
		{"google/corp/a/b/c/d/e/f/beta/foo.proto", 0},
		{"google/corp/a/b/c/foo.proto", 1},
	}

	for _, test := range tests {

		req, err := lint.NewProtoRequest(&descriptorpb.FileDescriptorProto{
			Name:           &test.path,
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
		})

		if err != nil {
			t.Errorf("Failed to create proto request because %v", err)
		}

		p, err := rule.Lint(req)

		if err != nil {
			t.Errorf("Lint() on file %q returned an error: %v", test.path, err)
		}

		if len(p) != test.numProblems {
			t.Errorf(
				"Lint() on file %q returned %d problems; want %d. Problems: %+v",
				test.path, len(p), test.numProblems, p,
			)
		}
	}
}
