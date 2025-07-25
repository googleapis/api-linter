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

package aip0165

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResponsePurgeSampleField(t *testing.T) {
	for _, test := range []struct {
		name     string
		RPC      string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "PurgeBooks", "repeated string purge_sample = 1;", nil},
		{"Missing", "PurgeBooks", "", testutils.Problems{{Message: "no `purge_sample`"}}},
		{"InvalidType", "PurgeBooks", "repeated int64 purge_sample = 1;", testutils.Problems{{Suggestion: "repeated string"}}},
		{"InvalidTypeNotRepeated", "PurgeBooks", "string purge_sample = 1;", testutils.Problems{{Suggestion: "repeated string"}}},
		{"IrrelevantMessage", "EnumerateBooks", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.RPC}}Response {
					{{.Field}}
					int32 purge_count = 2;
				}
			`, test)
			var d protoreflect.Descriptor = f.Messages().Get(0)
			if strings.HasPrefix(test.name, "InvalidType") {
				d = f.Messages().Get(0).Fields().Get(0)
			}
			if diff := test.problems.SetDescriptor(d).Diff(responsePurgeSampleField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
