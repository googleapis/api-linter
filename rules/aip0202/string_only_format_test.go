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

package aip0202

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestStringOnlyFormat(t *testing.T) {
	for _, tc := range []struct {
		name, Type, Format string
		problems           testutils.Problems
	}{
		{
			name:     "ValidSkipString",
			Type:     "string",
			Format:   "format: UUID4",
			problems: nil,
		},
		{
			name:     "ValidSkipNoFormat",
			Type:     "int32",
			problems: nil,
		},
		{
			name:     "InvalidNonString",
			Type:     "int32",
			Format:   "format: UUID4",
			problems: testutils.Problems{{Message: "string fields"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_info.proto";
			
			message CreateBookRequest {
				{{.Type}} foo = 1 [(google.api.field_info) = {
					{{.Format}}
				}];
			}
			`, tc)
			fd := file.Messages().ByName("CreateBookRequest").Fields().ByName("foo")
			problems := stringOnlyFormat.Lint(file)
			if diff := tc.problems.SetDescriptor(fd).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
