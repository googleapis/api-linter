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

package aip0151

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestResponseUnary(t *testing.T) {
	for _, test := range []struct {
		name     string
		Stream   string
		Response string
		problems testutils.Problems
	}{
		{"Vaild", "", "google.longrunning.Operation", nil},
		{"Invalid", "stream ", "google.longrunning.Operation", testutils.Problems{{Message: "unary, not stream"}}},
		{"Irrelevant", "stream ", "ReadBookResponse", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc ReadBook(ReadBookRequest) returns ({{.Stream}}{{.Response}});
				}
				message ReadBookRequest {}
				message ReadBookResponse {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(responseUnary.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
