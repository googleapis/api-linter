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
	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		methodName     string
		reqMessageName string
		problems       testutils.Problems
	}{
		{"Valid", "UpdateBook", "UpdateBookRequest", testutils.Problems{}},
		{"Invalid", "UpdateBook", "Book", testutils.Problems{{Suggestion: "UpdateBookRequest"}}},
		{"Irrelevant", "AcquireBook", "Book", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-134 Update method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage(test.reqMessageName), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
			)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the lint rule, and establish that it returns the expected problems.
			problems := requestMessageName.Lint(service.GetFile())
			if diff := test.problems.SetDescriptor(service.GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName     string
		MethodName   string
		RespTypeName string
		LRO          bool
		problems     testutils.Problems
	}{
		{"ValidResource", "UpdateBook", "Book", false, testutils.Problems{}},
		{"ValidLRO", "UpdateBook", "Book", true, testutils.Problems{}},
		{"Invalid", "UpdateBook", "UpdateBookResponse", false, testutils.Problems{{Suggestion: "Book"}}},
		{"InvalidLRO", "UpdateBook", "UpdateBookResponse", true, testutils.Problems{{Suggestion: "Book"}}},
		{"Irrelevant", "MutateBook", "MutateBookResponse", false, testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-134 Update method
			file := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns ({{ if .LRO }}google.longrunning.Operation{{ else }}{{ .RespTypeName }}{{ end }}) {
						{{ if .LRO -}}
						option (google.longrunning.operation_info) = {
							response_type: "{{.RespTypeName}}"
							metadata_type: "{{.MethodName}}Metadata"
						};
						{{ end -}}
					}
				}
				message {{.MethodName}}Request {}
				message {{.RespTypeName}} {}
			`, test)

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			problems := responseMessageName.Lint(file)
			method := file.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpMethod(t *testing.T) {
	// Set up POST and PATCH HTTP annotations.
	httpPost := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/{book.name=publishers/*/books/*}",
		},
	}
	httpPatch := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Patch{
			Patch: "/v1/{book.name=publishers/*/books/*}",
		},
	}

	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpRule   *annotations.HttpRule
		methodName string
		msg        string
	}{
		{"Valid", httpPatch, "UpdateBook", ""},
		{"Invalid", httpPost, "UpdateBook", "HTTP PATCH"},
		{"Irrelevant", httpPost, "AcquireBook", ""},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-134 Update method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("UpdateBookRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpMethod.Lint(service.GetFile())
			if test.msg == "" && len(problems) > 0 {
				t.Errorf("Got %v, expected no problems.", problems)
			} else if test.msg != "" && !strings.Contains(problems[0].Message, test.msg) {
				t.Errorf("Got %q, expected message containing %q", problems[0].Message, test.msg)
			}
		})
	}
}

func TestHttpBody(t *testing.T) {
	tests := []struct {
		testName   string
		body       string
		methodName string
		msg        string
	}{
		{"Valid", "book", "UpdateBook", ""},
		{"InvalidFoo", "foo", "UpdateBook", "HTTP body"},
		{"InvalidStar", "*", "UpdateBook", "HTTP body"},
		{"InvalidEmpty", "", "UpdateBook", "HTTP body"},
		{"Irrelevant", "*", "AcquireBook", ""},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Patch{
					Patch: "/v1/{book.name=publishers/*/books/*}",
				},
				Body: test.body,
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-134 Update method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("UpdateBookRequest"), false),
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

func TestSynonyms(t *testing.T) {
	tests := []struct {
		MethodName string
		problems   testutils.Problems
	}{
		{"UpdateBook", testutils.Problems{}},
		{"PatchBook", testutils.Problems{{Suggestion: "UpdateBook"}}},
		{"PutBook", testutils.Problems{{Suggestion: "UpdateBook"}}},
		{"SetBook", testutils.Problems{{Suggestion: "UpdateBook"}}},
	}
	for _, test := range tests {
		t.Run(test.MethodName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book);
				}
				message {{.MethodName}}Request {}
				message Book {}
			`, test)
			m := file.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(synonyms.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
