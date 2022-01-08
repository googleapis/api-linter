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

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestJavaOuterClassname(t *testing.T) {
	for _, test := range []struct {
		name       string
		filename   string
		statements []string
		problems   testutils.Problems
	}{
		{"Valid", "test.proto", []string{"package foo.v1;", `option java_outer_classname = "TestProto";`}, testutils.Problems{}},
		{"ValidWithPath", "foo/bar/test.proto", []string{"package foo.v1;", `option java_outer_classname = "TestProto";`}, testutils.Problems{}},
		{"ValidCompound", "appleorange.proto", []string{"package foo.v1;", `option java_outer_classname = "AppleOrangeProto";`}, testutils.Problems{}},
		{"ValidCompoundWithUnderscore", "apple_orange.proto", []string{"package foo.v1;", `option java_outer_classname = "AppleOrangeProto";`}, testutils.Problems{}},
		{"ValidCompoundWithPath", "foo/bar/appleorange.proto", []string{"package foo.v1;", `option java_outer_classname = "AppleOrangeProto";`}, testutils.Problems{}},
		{"InvalidWrong", "test.proto", []string{"package foo.v1;", `option java_outer_classname = "OtherProto";`}, testutils.Problems{{Message: `java_outer_classname = "TestProto"`}}},
		{"InvalidWrongWithPath", "foo/bar/test.proto", []string{"package foo.v1;", `option java_outer_classname = "OtherProto";`}, testutils.Problems{{Message: `java_outer_classname = "TestProto"`}}},
		{"InvalidMissing", "test.proto", []string{"package foo.v1;", ""}, testutils.Problems{{Message: `java_outer_classname = "TestProto"`}}},
		{"InvalidMissingWithPath", "foo/bar/test.proto", []string{"package foo.v1;"}, testutils.Problems{{Message: `java_outer_classname = "TestProto"`}}},
		{"Ignored", "test.proto", []string{"", ""}, testutils.Problems{}},
		{"IgnoredMaster", "test.proto", []string{"package foo.master;", ""}, testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			files := testutils.ParseProtoStrings(t, map[string]string{test.filename: strings.Join(test.statements, "\n")})
			f := files[test.filename]
			if diff := test.problems.SetDescriptor(f).Diff(javaOuterClassname.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
