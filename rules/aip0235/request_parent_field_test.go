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

package aip0235

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
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
		{"Valid", "", "BatchDeleteBooks", "string parent = 1;", "publishers/{p}/books/{b}", nil},
		{"Missing", "", "BatchDeleteBooks", "", "publishers/{p}/books/{b}", testutils.Problems{{Message: "no `parent`"}}},
		{"InvalidType", "", "BatchDeleteBooks", "int32 parent = 1;", "publishers/{p}/books/{b}", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantRPCName", "", "EnumerateBooks", "", "publishers/{p}/books/{b}", nil},
		{"IrrelevantNoParent", "", "BatchDeleteBooks", "", "books/{b}", nil},

		{"PackageValid", "package foo;", "BatchDeleteBooks", "string parent = 1;", "publishers/{p}/books/{b}", nil},
		{"PackageMissing", "package foo;", "BatchDeleteBooks", "", "publishers/{p}/books/{b}", testutils.Problems{{Message: "no `parent`"}}},
		{"PackageInvalidType", "package foo;", "BatchDeleteBooks", "int32 parent = 1;", "publishers/{p}/books/{b}", testutils.Problems{{Suggestion: "string"}}},
		{"PackageIrrelevantRPCName", "package foo;", "EnumerateBooks", "", "publishers/{p}/books/{b}", nil},
		{"PackageIrrelevantNoParent", "package foo;", "BatchDeleteBooks", "", "books/{b}", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				{{.Package}}
				import "google/api/resource.proto";
				import "google/protobuf/empty.proto";

				service Library {
					rpc {{.RPC}}({{.RPC}}Request) returns (google.protobuf.Empty);
				}

				message {{.RPC}}Request {
					{{.Field}}
				}

				message Book {
					option (google.api.resource) = {
						pattern: "{{.Pattern}}";
					};
				}
			`, test)
			var d protoreflect.Descriptor = f.Messages()[0]
			if test.name == "InvalidType" || test.name == "PackageInvalidType" {
				d = f.Messages()[0].Fields()[0]
			}
			if diff := test.problems.SetDescriptor(d).Diff(requestParentField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
