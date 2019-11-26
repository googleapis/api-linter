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

import "testing"

import "github.com/googleapis/api-linter/rules/internal/testutils"

func TestReservedWords(t *testing.T) {
	reservedWordsSet.Each(func(s string) {
		t.Run(s, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Book {
					string name = 1;
					string {{.N}} = 2;
				}
			`, struct{ N string }{N: s})
			field := f.GetMessageTypes()[0].GetFields()[1]
			want := testutils.Problems{{Message: field.GetName(), Descriptor: field}}
			if diff := want.Diff(reservedWords.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	})
}
