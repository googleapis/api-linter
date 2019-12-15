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
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestJavaMultipleFiles(t *testing.T) {
	for _, test := range []struct {
		name     string
		Package  string
		Opt      string
		problems testutils.Problems
	}{
		{"Valid", "package foo.v1;", "option java_multiple_files = true;", testutils.Problems{}},
		{"Invalid", "package foo.v1;", "", testutils.Problems{{Message: "java_multiple_files"}}},
		{"Ignored", "", "", testutils.Problems{}},
		{"IgnoredMaster", "package foo.master;", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				{{.Package}}
				{{.Opt}}
			`, test)
			if diff := test.problems.SetDescriptor(f).Diff(javaMultipleFiles.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
