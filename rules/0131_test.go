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

func TestGetRequestMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip131 {
		rpc GetFoo({{ .RequestName }}) returns (Foo);
	}

	message {{.RequestName }} {
		string name = 1;
	}

	message Foo {}`

	tests := []struct {
		RequestName  string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"GetFooRequest", 0, "", -1},
		{"GetFooReq", 1, "GetFooRequest", 5},
	}

	rule := checkGetRequestMessageName()

	for _, test := range tests {
		errPrefix := "AIP-131 Request Name"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{
				AdditionalProtoPaths: []string{testdatadir("api-common-protos")},
				Filename:             "test.proto",
				Template:             tmpl,
				Data:                 test,
			}), nil)
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

func TestGetURI(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	import "google/api/annotations.proto";

	service Aip131 {
		rpc GetBook(GetBookRequest) returns (Book) {
			option (google.api.http) = {
				{{ .Method }}: "{{ .URI }}"
			};
		}
	}

	message GetBookRequest {
		string name = 1;
	}

	message Book {}`

	tests := []struct {
		Method       string
		URI string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"get", "/v1/{name=publishers/*/books/*}", 0, "", -1},
		{"post", "/v1/{name=publishers/*/books/*}", 1, "", 7},
		{"get", "/v1/publishers/*/books/*", 1, "", 7},
		{"get", "/v1/publishers/{publisher_id}/books/{book_id}", 1, "", 7},
		{"get", "/v1/publishers/{publisher=*}/books/{book=*}", 1, "", 7},
		{"get", "/v1/{name=publishers/*/books/*}/somethingElse", 1, "", 7},
	}

	rule := checkGetURI()

	for _, test := range tests {
		errPrefix := "AIP-131 RPC URI"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{
				AdditionalProtoPaths: []string{testdatadir("api-common-protos")},
				Filename:             "test.proto",
				Template:             tmpl,
				Data:                 test,
			},
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

func TestGetBody(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	import "google/api/annotations.proto";

	service Aip131 {
		rpc GetBook(GetBookRequest) returns (Book) {
			option (google.api.http) = {
				get: "/v1/{name=publishers/*/books/*}"
				{{ .Body }}
			};
		}
	}

	message GetBookRequest {
		string name = 1;
	}

	message Book {}`

	tests := []struct {
		Body string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"", 0, "", -1},
		{"body: \"*\"", 1, "", 7},
	}

	rule := checkGetBody()

	for _, test := range tests {
		errPrefix := "AIP-131 RPC Body"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{
				AdditionalProtoPaths: []string{testdatadir("api-common-protos")},
				Filename:             "test.proto",
				Template:             tmpl,
				Data:                 test,
			},
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


func TestGetRequestMessageNameField(t *testing.T) {
	tmpl := `syntax = "proto3";

	service Aip131 {
		rpc GetFoo(GetFooRequest) returns (Foo);
	}

	message GetFooRequest {
		{{.NameFieldType}} {{.NameFieldName}} = 1;
	}

	message Foo {}`

	tests := []struct {
		NameFieldType string
		NameFieldName string
		problemCount  int
		startLine     int
	}{
		{"string", "name", 0, -1},
		{"string", "resource", 1, 7},
		{"bytes", "name", 1, 8},
	}

	rule := checkGetRequestMessageNameField()

	for _, test := range tests {
		errPrefix := "AIP-131 Request Name Field"
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

func TestGetRequestMessageUnknownFields(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	import "google/protobuf/field_mask.proto";

	service Aip131 {
		rpc GetFoo(GetFooRequest) returns (Foo);
	}

	message GetFooRequest {
		string name = 1;
		{{- range $i, $f := .ExtraFields }}
		{{ $f }}
		{{- end }}
	}

	message Foo {}

	enum FooView {
		FOO_VIEW_UNSPECIFIED = 0;
	}
	`

	tests := []struct {
		ExtraFields  []string
		problemCount int
		startLine    int
	}{
		{problemCount: 0, startLine: -1},
		{ExtraFields: []string{"FooView view = 2;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"google.protobuf.FieldMask read_mask = 2;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string application_id = 2;"}, problemCount: 1, startLine: 12},
		{ExtraFields: []string{
			"string application_id = 2;",
			"Foo foo = 3;",
		}, problemCount: 2, startLine: 12},
	}

	rule := checkGetRequestMessageUnknownFields()

	for _, test := range tests {
		errPrefix := "AIP-131 Request Unknown Fields"
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
			for _, problem := range resp {
				t.Errorf(problem.Message)
			}
		}

		if len(resp) > 0 {
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

func TestGetResponseMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip131 {
		rpc GetFoo(GetFooRequest) returns ({{ .ResponseName }});
	}

	message GetFooRequest {
		string name = 1;
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
		{"NotFoo", 1, "Foo", 5},
	}

	rule := checkGetResponseMessageName()

	for _, test := range tests {
		errPrefix := "AIP-131 Response Message Name"
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
