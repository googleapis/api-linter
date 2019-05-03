package corp

import (
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/lint"
	"testing"
)

func TestProtoPackageMustMatchFilePath(t *testing.T) {
	rule := protoPackageMustMatchFilePath()

	tests := []struct {
		path         string
		protoPackage string
		numProblems  int
	}{
		{"a/b/c/v1/foo.proto", "a.b.c.v1", 0},
		{"google/corp/a/b/c/d/e/f/beta/foo.proto", "google.corp.a.b.c.d.e.f.beta", 0},
		{"google/corp/a/b/c/foo.proto", "google.corp.a.b.c", 0},

		// should be "a.b.c"
		{"a/b/c/d.proto", "a.b.c.d", 1},
		{"a/b/c/d.proto", "a.b.d", 1},
		{"a/b/c/d.proto", "a.b", 1},
	}

	for _, test := range tests {

		req, err := lint.NewProtoRequest(&descriptorpb.FileDescriptorProto{
			Name:           &test.path,
			Package:        &test.protoPackage,
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
		})

		if err != nil {
			t.Errorf("Failed to create proto request because %v", err)
		}

		p, err := rule.Lint(req)

		if err != nil {
			t.Errorf("Lint() on file %q returned an error: %v", test.path, err)
		}

		if len(p.Problems) != test.numProblems {
			t.Errorf(
				"Lint() on file %q returned %d problems; want %d. Problems: %+v",
				test.path, len(p.Problems), test.numProblems, p.Problems,
			)
		}
	}
}
