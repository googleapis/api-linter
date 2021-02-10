// Copyright 2021 Google LLC
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

package aip0162

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRollbackRequestRevisionIDField(t *testing.T) {
	for _, test := range []struct {
		name     string
		RPC      string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "RollbackBook", "string revision_id = 1;", nil},
		{"Missing", "RollbackBook", "", testutils.Problems{{Message: "no `revision_id`"}}},
		{"InvalidType", "RollbackBook", "int32 revision_id = 1;", testutils.Problems{{Suggestion: "string"}}},
		{"IrrelevantRPCName", "EnumerateBooks", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.RPC}}Request {
					{{.Field}}
				}
			`, test)
			var d desc.Descriptor = f.GetMessageTypes()[0]
			if test.name == "InvalidType" {
				d = f.GetMessageTypes()[0].GetFields()[0]
			}
			if diff := test.problems.SetDescriptor(d).Diff(rollbackRequestRevisionIDField.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
