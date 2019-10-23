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

package aip0135

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
		{"Valid", "DeleteBook", "DeleteBookRequest", testutils.Problems{}},
		{"Invalid", "DeleteBook", "Book", testutils.Problems{{Suggestion: "DeleteBookRequest"}}},
		{"Irrelevant", "AcquireBook", "Book", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-135 Delete method
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
	tmpl := map[string]string{
		"sync": `
			package test;
			import "google/protobuf/empty.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespTypeName}});
			}
			message {{.MethodName}}Request {}
			{{ if (ne .RespTypeName "google.protobuf.Empty") }}message {{.RespTypeName}} {}{{ end }}
		`,
		"lro": `
			package test;
			import "google/longrunning/operations.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request)
				    returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "{{.RespTypeName}}"
						metadata_type: "{{.MethodName}}Metadata"
					};
				}
			}
			message {{.MethodName}}Request {}
		`,
	}

	// Set up the testing permutations.
	tests := []struct {
		testName     string
		tmpl         string
		MethodName   string
		RespTypeName string
		problems     testutils.Problems
	}{
		{"ValidEmpty", tmpl["sync"], "DeleteBook", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidResource", tmpl["sync"], "DeleteBook", "Book", testutils.Problems{}},
		{"ValidLROEmpty", tmpl["lro"], "DeleteBook", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidLROResource", tmpl["lro"], "DeleteBook", "Book", testutils.Problems{}},
		{"Invalid", tmpl["sync"], "DeleteBook", "DeleteBookResponse", testutils.Problems{{Suggestion: "google.protobuf.Empty"}}},
		{"InvalidLRO", tmpl["lro"], "DeleteBook", "DeleteBookResponse", testutils.Problems{{Suggestion: "google.protobuf.Empty"}}},
		{"Irrelevant", tmpl["sync"], "DestroyBook", "DestroyBookResponse", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a minimal service with a AIP-135 Delete method
			file := testutils.ParseProto3Tmpl(t, test.tmpl, test)

			// Run the lint rule, and establish that it returns the expected problems.
			method := file.GetServices()[0].GetMethods()[0]
			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpMethod(t *testing.T) {
	// Set up GET and DELETE HTTP annotations.
	httpGet := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{name=publishers/*/books/*}",
		},
	}
	httpDelete := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Delete{
			Delete: "/v1/{name=publishers/*/books/*}",
		},
	}

	// Set up testing permutations.
	tests := []struct {
		testName   string
		httpRule   *annotations.HttpRule
		methodName string
		msg        string
	}{
		{"Valid", httpDelete, "DeleteBook", ""},
		{"Invalid", httpGet, "DeleteBook", "HTTP DELETE"},
		{"Irrelevant", httpGet, "AcquireBook", ""},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-135 Delete method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("DeleteBookRequest"), false),
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
		{"Valid", "", "DeleteBook", ""},
		{"Invalid", "*", "DeleteBook", "HTTP body"},
		{"Irrelevant", "*", "AcquireBook", ""},
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

			// Create a minimal service with a AIP-135 Get method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("DeleteBookRequest"), false),
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
		{"Valid", "/v1/{name=publishers/*/books/*}", "DeleteBook", ""},
		{"InvalidVarName", "/v1/{book=publishers/*/books/*}", "DeleteBook", "`name` field"},
		{"NoVarName", "/v1/publishers/*/books/*", "DeleteBook", "`name` field"},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "AcquireBook", ""},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}
			httpRule := &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Delete{
					Delete: test.uri,
				},
			}
			if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-135 Delete method
			// (or with a different method, in the "Irrelevant" case).
			service, err := builder.NewService("Library").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage("DeleteBookRequest"), false),
				builder.RpcTypeMessage(builder.NewMessage("Book"), false),
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
