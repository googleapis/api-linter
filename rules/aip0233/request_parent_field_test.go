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

package aip0233

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestParentField(t *testing.T) {
	for _, test := range []struct {
		name     string
		RPC      string
		Field    string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "BatchCreateBooks", "string parent = 1;", "publishers/{p}/books", nil},
		{"Missing", "BatchCreateBooks", "", "publishers/{p}/books", testutils.Problems{{Message: "no `parent`"}}},
		{"InvalidType", "BatchCreateBooks", "int32 parent = 1;", "publishers/{p}/books", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantRPCName", "EnumerateBooks", "", "publishers/{p}/books", nil},
		{"IrrelevantNoParent", "BatchCreateBooks", "", "books", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
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
			if test.name == "InvalidType" {
				d = f.GetMessageTypes()[0].GetFields()[0]
			}
			if diff := test.problems.SetDescriptor(d).Diff(requestParentField.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
