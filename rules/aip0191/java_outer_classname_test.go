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

package aip0191

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestJavaOuterClassname(t *testing.T) {
	for _, test := range []struct {
		name       string
		statements []string
		problems   testutils.Problems
	}{
		{"Valid", []string{"package foo.v1;", `option java_outer_classname = "TestProto";`}, testutils.Problems{}},
		{"Invalid", []string{"package foo.v1;", ""}, testutils.Problems{{Message: `java_outer_classname = "TestProto"`}}},
		{"Ignored", []string{"", ""}, testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3String(t, strings.Join(test.statements, "\n"))
			if diff := test.problems.SetDescriptor(f).Diff(javaOuterClassname.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
