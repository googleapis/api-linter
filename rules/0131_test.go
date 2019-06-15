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

	"github.com/googleapis/api-linter/rules/testutil"
)

func TestRequestMessageName(t *testing.T) {
	tmpl := testutil.MustCreateTemplate(`
	syntax = "proto3";

	service Aip131 {
		rpc GetFoo({{ .RequestName }}) returns (Foo);
	}

	message {{.RequestName }} {
		string name = 1;
	}

	message Foo {}
	`)

	tests := []struct {
		RequestName  string
		problemCount int
		suggestion   string
		startLine    int
	}{
		{"GetFooRequest", 0, "", -1},
		{"GetFooReq", 1, "GetFooRequest", 5},
	}

	rule := checkRequestMessageName()

	for _, test := range tests {
		req := testutil.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := "AIP-131 Request Name"
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

func TestRequestMessageNameField(t *testing.T) {
	tmpl := testutil.MustCreateTemplate(`syntax = "proto3";

	service Aip131 {
		rpc GetFoo(GetFooRequest) returns (Foo);
	}

	message GetFooRequest {
		{{.NameFieldType}} {{.NameFieldName}} = 1;
	}

	message Foo {}
	`)

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

	rule := checkRequestMessageNameField()

	for _, test := range tests {
		req := testutil.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := "AIP-131 Request Name Field"
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

func TestRequestMessageUnknownFields(t *testing.T) {
	tmpl := testutil.MustCreateTemplate(`
	syntax = "proto3";

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
	`)

	tests := []struct {
		ExtraFields  []string
		problemCount int
		startLine    int
	}{
		{problemCount: 0, startLine: -1},
		{ExtraFields: []string{"string application_id = 2;"}, problemCount: 1, startLine: 10},
		{ExtraFields: []string{
			"string application_id = 2;",
			"Foo foo = 3;",
		}, problemCount: 2, startLine: 10},
	}

	rule := checkRequestMessageUnknownFields()

	for _, test := range tests {
		req := testutil.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := "AIP-131 Request Name Unknown Fields"
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
