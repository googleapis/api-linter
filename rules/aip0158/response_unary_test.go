// Copyright 2020 Google LLC
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

package aip0158

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseUnary(t *testing.T) {
	for _, test := range []struct {
		name          string
		Stream        string
		Response      string
		NextPageToken string
		problems      testutils.Problems
	}{
		{"Vaild", "", "ListFoosResponse", "next_page_token", nil},
		{"Invalid", "stream ", "ListFoosResponse", "next_page_token", testutils.Problems{{Message: "unary, not stream"}}},
		{"Irrelevant", "stream ", "FrobFoosResponse", "not_paginated_response", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Foos {
					rpc ListFoos(ListFoosRequest) returns ({{.Stream}}{{.Response}});
				}
				message ListFoosRequest {
					int32 page_size = 1;
					string page_token = 2;
				}
				message {{.Response}} {
					string {{.NextPageToken}} = 1;
				}
			`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(responseUnary.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
