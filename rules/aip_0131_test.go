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

package rules

import (
	"testing"

	"github.com/googleapis/api-linter/rules/testutil"
)

func TestAIP0131(t *testing.T) {
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
		RequestName string
		numProblem  int
		suggestion  string
		startLine   int
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

		if got, want := len(resp), test.numProblem; got != want {
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
