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

package rules

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/testutil"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestGetHTTPRulesEmpty(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	import "google/api/annotations.proto";

	service Library {
		rpc GetBook(GetBookRequest) returns (Book) {
			{{ .Option }}
		}
	}

	message GetBookRequest {
		string name = 1;
	}

	message Book {}`

	tests := []struct {
		Option    string
		ruleCount int
	}{
		{"", 0},
		{"option (google.api.http) = { get: \"/v1/publishers/*/books/*\" };", 1},
		{"option (google.api.http) = { get: \"/v1/publishers/*/books/*\" additional_bindings: { get: \"/v1/authors/*/books/*\" } };", 2},
	}

	for _, test := range tests {
		fdp := testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{
				AdditionalProtoPaths: []string{testdatadir("api-common-protos")},
				Filename:             "test.proto",
				Template:             tmpl,
				Data:                 test,
			},
		)
		reg, err := lint.MakeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fdp})
		if err != nil {
			t.Fatalf("getHTTPRules: Failed to initialize registry.")
		}
		fd, err := protodesc.NewFile(fdp, reg)
		if err != nil {
			t.Fatalf("getHTTPRules: Failed to initialize proto")
		}
		method := fd.Services().Get(0).Methods().Get(0)
		httpRules := getHTTPRules(method)
		if got, expected := len(httpRules), test.ruleCount; got != expected {
			t.Errorf("getHTTPRules: Expected %d rules, but got %d.", expected, got)
		}
	}
}
