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

package aip0133

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName      string
		resourceField string
		body          string
		methodName    string
		msg           string
	}{
		{"Valid", "book", "book", "CreateBook", ""},
		{"Valid", "textbook", "textbook", "CreateBook", ""},
		{"Valid", "", "book", "CreateBook", ""}, // valid for http body rule check, but it will fail under resource fail rule check
		{"Invalid_BodyMissing", "book", "", "CreateBook", "Post methods should have an HTTP body"},
		{"Invalid_BodyMismatch", "book", "abook", "CreateBook", "The content of body \"abook\" must map to the resource field \"book\" in the request message"},
		{"Irrelevant", "book", "book", "CreateBook", ""},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{parent=publishers/*/books/*}",
				},
				Body: test.body,
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			book, err := builder.NewMessage("Book").Build()

			if err != nil {
				t.Fatalf("Failed to create resource 'Book' message for test.")
			}

			// Create a minimal service with a AIP-133 Create method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("CreateBookRequest").AddField(builder.NewField(test.resourceField, builder.FieldTypeImportedMessage(book))), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpBody.Lint(service.GetFile())
			if test.msg == "" && len(problems) > 0 {
				t.Errorf("Got %v, expected no problems.", problems)
			} else if test.msg != "" && !strings.Contains(problems[0].Message, test.msg) {
				t.Errorf("Got %q, expected message containing %q", problems[0].Message, test.msg)
			}
		})
	}
}
