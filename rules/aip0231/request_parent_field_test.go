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

package aip0231

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestParentField(t *testing.T) {
	for _, test := range []struct {
		name     string
		Package  string
		RPC      string
		Field    string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "", "BatchGetBooks", "string parent = 1;", "publishers/{p}/books/{b}", nil},
		{"Missing", "", "BatchGetBooks", "", "publishers/{p}/books/{b}", testutils.Problems{{Message: "no `parent`"}}},
		{"InvalidType", "", "BatchGetBooks", "int32 parent = 1;", "publishers/{p}/books/{b}", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantRPCName", "", "EnumerateBooks", "", "publishers/{p}/books/{b}", nil},
		{"IrrelevantNoParent", "", "BatchGetBooks", "", "books/{b}", nil},

		{"PackageValid", "package foo;", "BatchGetBooks", "string parent = 1;", "publishers/{p}/books/{b}", nil},
		{"PackageMissing", "package foo;", "BatchGetBooks", "", "publishers/{p}/books/{b}", testutils.Problems{{Message: "no `parent`"}}},
		{"PackageInvalidType", "package foo;", "BatchGetBooks", "int32 parent = 1;", "publishers/{p}/books/{b}", testutils.Problems{{Suggestion: "string"}}},
		{"PackageIrrelevantRPCName", "package foo;", "EnumerateBooks", "", "publishers/{p}/books/{b}", nil},
		{"PackageIrrelevantNoParent", "package foo;", "BatchGetBooks", "", "books/{b}", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				{{.Package}}
				import "google/api/resource.proto";

				service Library {
					rpc {{.RPC}}({{.RPC}}Request) returns ({{.RPC}}Response);
				}

				message {{.RPC}}Request {
					{{.Field}}
					repeated string names = 2;
				}

				message {{.RPC}}Response {
					repeated Book books = 1;
				}

				message Book {
					option (google.api.resource) = {
						pattern: "{{.Pattern}}";
					};
					string name = 1;
				}
			`, test)
			var d desc.Descriptor = f.GetMessageTypes()[0]
			if test.name == "InvalidType" || test.name == "PackageInvalidType" {
				d = f.GetMessageTypes()[0].GetFields()[0]
			}
			if diff := test.problems.SetDescriptor(d).Diff(requestParentField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
