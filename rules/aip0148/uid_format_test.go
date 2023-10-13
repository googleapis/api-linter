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

package aip0148

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestUidFormat(t *testing.T) {
	for _, test := range []struct {
		name, FieldName, Annotation string
		problems                    testutils.Problems
	}{
		{
			name:       "ValidUidFormat",
			FieldName:  "uid",
			Annotation: "[(google.api.field_info).format = UUID4]",
		},
		{
			name:      "SkipNonUid",
			FieldName: "other",
		},
		{
			name:      "InvalidMissingFormat",
			FieldName: "uid",
			problems:  testutils.Problems{{Message: "format = UUID4"}},
		},
		{
			name:       "InvalidWrongFormat",
			FieldName:  "uid",
			Annotation: "[(google.api.field_info).format = IPV4]",
			problems:   testutils.Problems{{Message: "format = UUID4"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_info.proto";

				message Person {
					string {{.FieldName}} = 2 {{.Annotation}};
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(uidFormat.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
