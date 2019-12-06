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

package aip0214

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestTtlType(t *testing.T) {
	for _, test := range []struct {
		name     string
		Type     string
		problems testutils.Problems
	}{
		{"Valid", "google.protobuf.Duration", testutils.Problems{}},
		{"Invalid", "int32", testutils.Problems{{Suggestion: "google.protobuf.Duration"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/protobuf/duration.proto";
			message Book {
				string name = 1;
				{{.Type}} ttl = 2;
			}
		`, test)
		field := f.GetMessageTypes()[0].GetFields()[1]
		if diff := test.problems.SetDescriptor(field).Diff(ttlType.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	}
}
