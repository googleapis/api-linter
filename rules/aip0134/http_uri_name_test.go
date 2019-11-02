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

package aip0134

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestHttpNameField(t *testing.T) {
	tests := []struct {
		testName   string
		uri        string
		methodName string
		msg        string
	}{
		{"Valid", "/v1/{big_book.name=publishers/*/books/*}",
			"UpdateBigBook", ""},
		{"InvalidNoUnderscore", "/v1/{bigbook.name=publishers/*/books/*}",
			"UpdateBigBook", "`big_book.name` field"},
		{"InvalidVarNameBook", "/v1/{big_book=publishers/*/books/*}",
			"UpdateBigBook", "`big_book.name` field"},
		{"InvalidVarNameName", "/v1/{name=publishers/*/books/*}",
			"UpdateBigBook", "`big_book.name` field"},
		{"InvalidVarNameReversed", "/v1/{name.big_book=publishers/*/books/*}",
			"UpdateBigBook", "`big_book.name` field"},
		{"NoVarName", "/v1/publishers/*/books/*",
			"UpdateBigBook", "`big_book.name` field"},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}",
			"AcquireBigBook", ""},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Patch{
					Patch: test.uri,
				},
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-134 Update method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("UpdateBigBookRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("BigBook"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpNameField.Lint(service.GetFile())
			if test.msg == "" && len(problems) > 0 {
				t.Errorf("Got %v, expected no problems.", problems)
			} else if test.msg != "" && len(problems) == 0 {
				t.Errorf("Got no problems, expected 1.")
			} else if test.msg != "" && !strings.Contains(problems[0].Message, test.msg) {
				t.Errorf("Got %q, expected message containing %q", problems[0].Message, test.msg)
			}
		})
	}
}
