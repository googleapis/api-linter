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

// TODO: Use template for testing.
func TestHttpVerb(t *testing.T) {
	// Set up GET and POST HTTP annotations.
	httpGet := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{parent=publishers/*}/books:batchGet",
		},
	}
	httpPost := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/{name=publishers/*/books/*}",
		},
	}

	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpRule   *annotations.HttpRule
		methodName string
		problems   testutils.Problems
	}{
		{"Valid", httpGet, "BatchGetBooks", nil},
		{"Invalid", httpPost, "BatchGetBooks", testutils.Problems{{Message: "Batch Get methods must use the HTTP GET verb."}}},
		{"Irrelevant", httpPost, "AcquireBook", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
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
			method := service.GetMethods()[0]
			problems := httpVerb.Lint(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
