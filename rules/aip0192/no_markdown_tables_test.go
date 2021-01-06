// Copyright 2021 Google LLC
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

func TestNoMarkdownTables(t *testing.T) {
	problem := testutils.Problems{{Message: "Markdown tables"}}
	for _, test := range []struct {
		name     string
		Comment  string
		problems testutils.Problems
	}{
		{"Valid", "foo bar baz", nil},
		{"ValidMidline", "foo - bar", nil},
		{"ValidPipe", "foo | bar", nil},
		{"WeirdButLegal", "|-|-|", nil},
		{"TableSeparator", "--- | ---", problem},
		{"BoundedTableSeparator", "| --- | --- |", problem},
		{"MoreColumnTable", "--- | --- | --- | --- | ---", problem},
		{"BiggerColumnTable", "------------ | -----", problem},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			  // {{.Comment}}
			  message Foo {}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(noMarkdownTables.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
