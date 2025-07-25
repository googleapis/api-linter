// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNoPrepositions(t *testing.T) {
	tests := []struct {
		Name     string
		problems testutils.Problems
	}{
		{"author", testutils.Problems{}},
		{"written_by", testutils.Problems{{Message: "by"}}},
		{"move_toward_shelf_at_front", testutils.Problems{{Message: "toward"}, {Message: "at"}}},
		{"order_by", testutils.Problems{}},
		{"group_by", testutils.Problems{}},
		{"hour_of_day", testutils.Problems{}},
		{"day_of_week", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message Book {
					string {{.Name}} = 1;
				}
			`, test)
			field := file.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(noPrepositions.Lint(file)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
