// Copyright 2026 Google LLC
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

package aip0190

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestMethodCase(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		problems   testutils.Problems
	}{
		{"Valid", "GetBook", testutils.Problems{}},
		{"InvalidCustomerID", "CustomerID", testutils.Problems{{Message: `Method name "CustomerID" must use UpperCamelCase.`}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			service Library {
				rpc {{ .MethodName }}(GetBookRequest) returns (Book) {}
			}
			message GetBookRequest {}
			message Book {}
		`, test)
			m := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(m).Diff(methodCase.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
