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

func TestListRequestMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip132 {
		rpc ListFoos({{ .RequestName }}) returns (ListFoosResponse);
	}

	message {{.RequestName }} {
		string parent = 1;
		int32 page_size = 2;
		string page_token = 3;
	}

	message ListFoosResponse {}
	`

	tests := []struct {
		RequestName  string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"ListFoosRequest", 0, "", -1},
		{"ListFoosReq", 1, "ListFoosRequest", 5},
	}

	rule := checkListRequestMessageName()

	for _, test := range tests {
		errPrefix := "AIP-132 Request Name"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Template: tmpl, Data: test},
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

func TestListRequestMessageParentField(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip132 {
		rpc ListFoos(ListFoosRequest) returns (ListFoosResponse);
	}

	message ListFoosRequest {
		{{ .ParentFieldType }} {{ .ParentFieldName }} = 1;
		int32 page_size = 2;
		string page_token = 3;
	}

	message ListFoosResponse {}
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

	rule := checkListRequestMessageParentField()

	for _, test := range tests {
		errPrefix := "AIP-132 Request Parent Field"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Template: tmpl, Data: test},
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

func TestListRequestMessageUnknownFields(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip132 {
		rpc ListFoos(ListFoosRequest) returns (ListFoosResponse);
	}

	message ListFoosRequest {
		string parent = 1;
		int32 page_size = 2;
		string page_token = 3;
		{{- range $i, $f := .ExtraFields }}
		{{ $f }}
		{{- end }}
	}

	message ListFoosResponse {}

	message Foo {}
	`

	tests := []struct {
		ExtraFields  []string
		problemCount int
		startLine    int
	}{
		{problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string filter = 4;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string order_by = 4;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{
			"string filter = 4;",
			"string order_by = 5;",
		}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string group_by = 4;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"bool show_deleted = 4;"}, problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string application_id = 4;"}, problemCount: 1, startLine: 12},
		{ExtraFields: []string{
			"string application_id = 4;",
			"Foo foo = 5;",
		}, problemCount: 2, startLine: 12},
	}

	rule := checkListRequestMessageUnknownFields()

	for _, test := range tests {
		errPrefix := "AIP-132 Request Unknown Fields"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Template: tmpl, Data: test},
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

func TestListResponseMessageName(t *testing.T) {
	tmpl := `
	syntax = "proto3";

	service Aip132 {
		rpc ListFoos(ListFoosRequest) returns ({{ .ResponseName }});
	}

	message ListFoosRequest {
		string parent = 1;
		int32 page_size = 2;
		string page_token = 3;
	}

	message {{ .ResponseName }} {}
	`

	tests := []struct {
		ResponseName string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"ListFoosResponse", 0, "", -1},
		{"NotListFoosResponse", 1, "ListFoosResponse", 5},
	}

	rule := checkListResponseMessageName()

	for _, test := range tests {
		errPrefix := "AIP-132 Response Message Name"
		req, err := lint.NewProtoRequest(testutil.MustCreateFileDescriptorProto(
			t,
			testutil.FileDescriptorSpec{Template: tmpl, Data: test},
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
