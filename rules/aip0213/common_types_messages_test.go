// Copyright 2023 Google LLC
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

package aip0213

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestCommonTypesMessages(t *testing.T) {
	for _, test := range []struct {
		name     string
		Field1   string
		Field2   string
		Field3   string
		problems testutils.Problems
	}{
		{"ValidCommonType", "google.protobuf.Duration duration", "", "", nil},
		{"ValidOtherType", "string bar", "", "", nil},
		{"ValidMessageDoesNotHaveAllFieldsFromCommonType", "string red", "string green", "", nil},
		{"MessageHasAllFieldsFromCommonType", "string red", "string green", "string blue", testutils.Problems{{Message: "common type"}}},
		{"MessageHasReorderedFieldsFromCommonType", "string green", "string blue", "string red", testutils.Problems{{Message: "common type"}}},
		{"MessageHasAdditionalFieldsOutsideCommonType", "double latitude", "double longitude", "double altitude", testutils.Problems{{Message: "common type"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/duration.proto";

				message Foo {
					{{.Field1}} = 1;

					{{if .Field2}}
					{{.Field2}} = 2;
					{{end}}

					{{if .Field3}}
					{{.Field3}} = 3;
					{{end}}
				}
			`, test)
			message := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(message).Diff(commonTypesMessages.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
