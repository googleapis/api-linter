// Copyright 2025 Google LLC
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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestTimeOffsetType(t *testing.T) {
	for _, test := range []struct {
		name      string
		FieldName string
		FieldType string
		problems  testutils.Problems
	}{
		{
			name:      "Valid",
			FieldName: "start_time_offset",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
		{
			name:      "InvalidNotDuration",
			FieldName: "start_time_offset",
			FieldType: "string",
			problems:  testutils.Problems{{Suggestion: "google.protobuf.Duration"}},
		},
		{
			name:      "InvalidNotDurationInt",
			FieldName: "event_time_offset",
			FieldType: "int32",
			problems:  testutils.Problems{{Suggestion: "google.protobuf.Duration"}},
		},
		{
			name:      "IrrelevantNoTimeOffsetSuffix",
			FieldName: "start_offset",
			FieldType: "float",
			problems:  nil,
		},
		{
			name:      "IrrelevantNoSuffix",
			FieldName: "start",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/duration.proto";
				message AudioSegment {
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(timeOffsetType.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
