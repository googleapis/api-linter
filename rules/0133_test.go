// Copyright 2019 Google LLC
//
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
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
)

func TestCreateRequestMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip133 {
		rpc CreateFoo({{ .RequestName }}) returns (Foo);
	}

	message {{.RequestName }} {
		string parent = 1;
		Foo foo = 2;
	}

	message Foo {}
	`

	tests := []struct {
		RequestName  string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"CreateFooRequest", 0, "", -1},
		{"CreateFooReq", 1, "CreateFooRequest", 5},
	}

	rule := checkCreateRequestMessageName()

	for _, test := range tests {
		errPrefix := "AIP-133 Request Name"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Filename: "test.proto", Template: tmpl, Data: test},
		), nil)
		if err != nil {
			t.Errorf("%s: lint.NewProtoRequest returned error %v", errPrefix, err)
		}

		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp), test.problemCount; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp) > 0 {
			if got, want := resp[0].Suggestion, test.suggestion; got != want {
				t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, got, want)
			}
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

func TestCreateRequestMessageParentField(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip133 {
		rpc CreateFoo(CreateFooRequest) returns (Foo);
	}

	message CreateFooRequest {
		{{ .ParentFieldType }} {{ .ParentFieldName }} = 1;
		Foo foo = 2;
	}

	message Foo {}
	`

	tests := []struct {
		ParentFieldType string
		ParentFieldName string
		problemCount    int
		startLine       int
	}{
		{"string", "parent", 0, -1},
		{"string", "resource", 1, 8},
		{"bytes", "parent", 1, 9},
	}

	rule := checkCreateRequestMessageParentField()

	for _, test := range tests {
		errPrefix := "AIP-133 Request Parent Field"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Filename: "test.proto", Template: tmpl, Data: test},
		), nil)
		if err != nil {
			t.Errorf("%s: lint.NewProtoRequest returned error %v", errPrefix, err)
		}

		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp), test.problemCount; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp) > 0 {
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

func TestCreateRequestMessageUnknownFields(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip133 {
		rpc CreateFoo(CreateFooRequest) returns (Foo);
	}

	message CreateFooRequest {
		string parent = 1;
		Foo foo = 2;
		{{- range $i, $f := .ExtraFields }}
		{{ $f }}
		{{- end }}
	}

	message Foo {}
	`

	tests := []struct {
		ExtraFields  []string
		problemCount int
		startLine    int
	}{
		{problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string foo_id = 3;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string request_id = 3;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string application_id = 4;"}, problemCount: 1, startLine: 11},
		{ExtraFields: []string{
			"string foo_id = 3;",
			"string application_id = 4;",
		}, problemCount: 1, startLine: 12},
	}

	rule := checkCreateRequestMessageUnknownFields()

	for _, test := range tests {
		errPrefix := "AIP-133 Request Unknown Fields"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Filename: "test.proto", Template: tmpl, Data: test},
		), nil)
		if err != nil {
			t.Errorf("%s: lint.NewProtoRequest returned error %v", errPrefix, err)
		}

		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp), test.problemCount; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp) > 0 {
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

func TestCreateResponseMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip133 {
		rpc CreateFoo(CreateFooRequest) returns ({{ .ResponseName }});
	}

	message CreateFooRequest {
		string parent = 1;
		{{ .ResponseName }} foo = 2;
	}

	message {{ .ResponseName }} {}
	`

	tests := []struct {
		ResponseName string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"Foo", 0, "", -1},
		{"Operation", 0, "", -1},
		{"CreateFooResponse", 1, "Foo", 5},
	}

	rule := checkCreateResponseMessageName()

	for _, test := range tests {
		errPrefix := "AIP-133 Response Message Name"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Filename: "test.proto", Template: tmpl, Data: test},
		), nil)
		if err != nil {
			t.Errorf("%s: lint.NewProtoRequest returned error %v", errPrefix, err)
		}

		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp), test.problemCount; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp) > 0 {
			if got, want := resp[0].Suggestion, test.suggestion; got != want {
				t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, got, want)
			}
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}
