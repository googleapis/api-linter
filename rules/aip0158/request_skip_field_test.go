// Copyright 2021 Google LLC
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

package aip0158

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRequestSkipField(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		Fields      string
		problems    testutils.Problems
	}{
		{"ValidPresent", "FooRequest", `int32 page_size = 1; int32 skip = 2;`, nil},
		{"ValidAbsent", "FooRequest", `int32 page_size = 1;`, nil},
		{"InvalidType", "FooRequest", `int32 page_size = 1; string skip = 2;`, testutils.Problems{{Suggestion: "int32"}}},
		{"InvalidRepeated", "FooRequest", `int32 page_size = 1; repeated int32 skip = 2;`, testutils.Problems{{Suggestion: "int32"}}},
		{"IrrelevantMessageName", "Foo", `int32 page_size = 1; string skip = 2;`, nil},
		{"IrrelevantMessageFields", "FooRequest", `string skip = 1;`, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.Fields}}
				}
			`, test)
			var d protoreflect.Descriptor = nil
			if len(test.problems) > 0 {
				d = f.Messages()[0].Fields()[1]
			}
			problems := requestSkipField.Lint(f)
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
