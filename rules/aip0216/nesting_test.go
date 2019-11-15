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

package aip0216

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNesting(t *testing.T) {
	t.Run("Nested", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			message Foo {
				enum State {
					STATE_UNSPECIFIED = 0;
				}
			}
		`)
		if diff := (testutils.Problems{}).Diff(nesting.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	})
	t.Run("NotNested", func(t *testing.T) {
		tests := []struct {
			name        string
			MessageName string
			EnumName    string
			PackageStmt string
			problems    testutils.Problems
		}{
			{"ShouldNest", "Foo", "FooState", "package test;", testutils.Problems{{Message: "Nest"}}},
			{"ShouldNestNoPkg", "Foo", "FooState", "", testutils.Problems{{Message: "Nest"}}},
			{"ValidNoFoo", "Bar", "FooState", "package test;", testutils.Problems{}},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					{{.PackageStmt}}
					message {{.MessageName}} {}
					enum {{.EnumName}} {
						UNSPECIFIED = 0;
					}
				`, test)
				e := f.GetEnumTypes()[0]
				if diff := test.problems.SetDescriptor(e).Diff(nesting.Lint(f)); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})
}
