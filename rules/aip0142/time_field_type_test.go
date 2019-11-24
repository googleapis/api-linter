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

package aip0142

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestFieldType(t *testing.T) {
	tests := []struct {
		testName  string
		FieldType string
		FieldName string
		problems  testutils.Problems
	}{
		{"Valid", "google.protobuf.Timestamp", "create_time", testutils.Problems{}},
		{"Invalid", "int32", "created_ms", testutils.Problems{{Suggestion: "google.protobuf.Timestamp"}}},
		{"Irrelevant", "int32", "foo", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/timestamp.proto";
				message Book {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			problems := fieldType.Lint(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
