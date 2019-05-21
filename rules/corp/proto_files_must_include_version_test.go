// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package corp

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
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
