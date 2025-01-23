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

func TestCommonTypesFields(t *testing.T) {
	for _, test := range []struct {
		name     string
		Field    string
		problems testutils.Problems
	}{
		{"ValidCommonType", "google.protobuf.Duration duration", nil},
		{"ValidFieldEndingInCommonTerm", "google.protobuf.Duration max_duration", nil},
		{"ValidFieldEndingInMultipleCommonTerms", "google.type.TimeOfDay end_time_of_day", nil},
		{"ValidOtherType", "string bar", nil},
		{"FieldDoesNotUseMessageType", "map<int32, int32> duration", testutils.Problems{{Message: "common type"}}},
		{"FieldEndingInCommonTermDoesNotUseMessageType", "map<int32, int32> max_duration", testutils.Problems{{Message: "common type"}}},
		{"FieldEndingInMultipleCommonTermsDoesNotUseMessageType", "map<int32, int32> end_time_of_day", testutils.Problems{{Message: "common type"}}},
		{"FieldDoesNotUseCommonType", "int32 duration", testutils.Problems{{Message: "common type"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/duration.proto";
				import "google/type/timeofday.proto";

				message Foo {
					{{.Field}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(commonTypesFields.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
