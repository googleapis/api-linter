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

package aip0192

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestTrademarkedNames(t *testing.T) {
	for _, test := range []struct {
		Token    string
		problems testutils.Problems
	}{
		{"GitHub", testutils.Problems{}},
		{"Bigtable", testutils.Problems{}},
		{"Pub/Sub", testutils.Problems{}},
		{"Git Hub", testutils.Problems{{Message: "GitHub"}}},
		{"Git  Hub", testutils.Problems{{Message: "GitHub"}}},
		{"Git\n// Hub", testutils.Problems{{Message: "GitHub"}}},
		{"Git \n// Hub", testutils.Problems{{Message: "GitHub"}}},
		{"Git Hub Github", testutils.Problems{{Message: "GitHub"}, {Message: "GitHub"}}},
		{"G-Suite", testutils.Problems{{Message: "G Suite"}}},
		{"PubSub", testutils.Problems{{Message: "Pub/Sub"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			// This is a comment that says {{.Token}}.
			message Foo {}
		`, test)
		m := f.Messages().Get(0)
		if diff := test.problems.SetDescriptor(m).Diff(trademarkedNames.Lint(f)); diff != "" {
			t.Error(diff)
		}
	}
}
