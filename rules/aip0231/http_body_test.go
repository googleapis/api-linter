// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0231

import (
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName   string
		body       string
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", "", "BatchGetBooks", nil},
		{"Invalid", "*", "BatchGetBooks", testutils.Problems{{Message: "Batch Get methods should not have an HTTP body."}}},
		{"Irrelevant", "*", "AcquireBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{
					Get: "/v1/{name=publishers/*/books/*}",
				},
				Body: test.body,
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-231 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("BookService").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BatchGetBooksResponse"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpBody.Lint(service.GetFile())
			if diff := test.problems.SetDescriptor(service.GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
