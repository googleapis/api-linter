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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestDurationOffsetComment(t *testing.T) {
	for _, test := range []struct {
		name      string
		Comment   string
		FieldName string
		FieldType string
		problems  testutils.Problems
	}{
		{
			name:      "ValidWithComment",
			Comment:   "// The duration relative to the start of the stream.",
			FieldName: "start_offset",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
		{
			name:      "InvalidNoComment",
			Comment:   "",
			FieldName: "start_offset",
			FieldType: "google.protobuf.Duration",
			problems:  testutils.Problems{{Message: "must include a clear comment explaining the relative start point."}},
		},
		{
			name:      "ValidWithRespectComment",
			Comment:   "// The duration in respect to the start.",
			FieldName: "end_offset",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
		{
			name:      "ValidOfTheComment",
			Comment:   "// The duration of the event offset from start.",
			FieldName: "event_offset",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
		{
			name:      "InvalidInadequateComment",
			Comment:   "// This is just a comment.",
			FieldName: "another_offset",
			FieldType: "google.protobuf.Duration",
			problems:  testutils.Problems{{Message: "must include a clear comment explaining the relative start point."}},
		},
		{
			name:      "IrrelevantNoOffsetSuffix",
			Comment:   "",
			FieldName: "start",
			FieldType: "google.protobuf.Duration",
			problems:  nil,
		},
		{
			name:      "IrrelevantNotDuration",
			Comment:   "",
			FieldName: "map_offset",
			FieldType: "string",
			problems:  nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/duration.proto";
				message AudioSegment {
					{{.Comment}}
					{{.FieldType}} {{.FieldName}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(durationOffsetComment.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
