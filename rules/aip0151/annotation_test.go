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

package aip0151

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestAnnotationExistsValid(t *testing.T) {
	f := testutils.ParseProto3String(t, `
		import "google/longrunning/operations.proto";
		service Library {
			rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "WriteBookResponse"
					metadata_type: "WriteBookMetadata"
				};
			}
		}
		message WriteBookRequest {}
	`)
	if diff := (testutils.Problems{}).Diff(lroAnnotationExists.Lint(f)); diff != "" {
		t.Errorf(diff)
	}
}

func TestAnnotationExistsInvalid(t *testing.T) {
	f := testutils.ParseProto3String(t, `
		import "google/longrunning/operations.proto";
		service Library {
			rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation);
		}
		message WriteBookRequest {}
	`)
	want := testutils.Problems{{
		Descriptor: f.GetServices()[0].GetMethods()[0],
		Message:    "operation_info annotation",
	}}
	if diff := want.Diff(lroAnnotationExists.Lint(f)); diff != "" {
		t.Errorf(diff)
	}
}

func TestLROResponse(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		Response   string
		problems   testutils.Problems
	}{
		{"Valid", "WriteBook", "WriteBookResponse", testutils.Problems{}},
		{"InvalidEmptyString", "WriteBook", "", testutils.Problems{{Message: "must set the response type"}}},
		{"InvalidGPEmpty", "WriteBook", "google.protobuf.Empty", testutils.Problems{{Message: "Empty"}}},
		{"ValidGPEmptyDelete", "DeleteBook", "google.protobuf.Empty", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.Response}}"
							metadata_type: "{{.MethodName}}Metadata"
						};
					}
				}
				message {{.MethodName}}Request {}
			`, test)
			problems := lroResponse.Lint(f)
			d := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestLROMetadata(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		Metadata   string
		problems   testutils.Problems
	}{
		{"Valid", "WriteBook", "WriteBookMetadata", testutils.Problems{}},
		{"InvalidEmptyString", "WriteBook", "", testutils.Problems{{Message: "must set the metadata type"}}},
		{"InvalidGPEmpty", "WriteBook", "google.protobuf.Empty", testutils.Problems{{Message: "Empty"}}},
		{"InvalidGPEmptyDelete", "DeleteBook", "google.protobuf.Empty", testutils.Problems{{Message: "Empty"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.MethodName}}Response"
							metadata_type: "{{.Metadata}}"
						};
					}
				}
				message {{.MethodName}}Request {}
			`, test)
			problems := lroMetadata.Lint(f)
			d := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestDefinedInFile(t *testing.T) {
	tests := []struct {
		testName     string
		Package      string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "", "WriteBookResponse", testutils.Problems{}},
		{"ValidPkg", "package test;", "WriteBookResponse", testutils.Problems{}},
		{"InvalidTypo", "", "WriteBookRepsonse", testutils.Problems{{Message: "should be defined"}}},
		{"InvalidTypoPkg", "package test;", "WriteBookRepsonse", testutils.Problems{{Message: "should be defined"}}},
		{"ValidExternal", "", "google.protobuf.Empty", testutils.Problems{}},
		{"ValidExternalPkg", "package test;", "google.protobuf.Empty", testutils.Problems{}},
	}
	for _, test := range tests {
		f := testutils.ParseProto3Tmpl(t, `
			{{.Package}}
			import "google/longrunning/operations.proto";
			service Library {
				rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
					option (google.longrunning.operation_info) = {
						response_type: "{{.ResponseType}}"
						metadata_type: "OperationMetadata"
					};
				}
			}
			message WriteBookRequest {}
			message WriteBookResponse {}
			message OperationMetadata {}
		`, test)
		problems := lroDefinedInFile.Lint(f)
		m := f.GetServices()[0].GetMethods()[0]
		if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
			t.Errorf(diff)
		}
	}
}
