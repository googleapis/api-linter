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

package aip0148

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestFieldBehavior(t *testing.T) {
	const messageOptsResource = `option (google.api.resource).type = "library.googleapis.com/Book";`
	const fieldOptsOutputOnly = `[(google.api.field_behavior) = OUTPUT_ONLY]`
	missingOutputOnly := testutils.Problems{{Message: "OUTPUT_ONLY"}}

	for _, test := range []struct {
		name        string
		MessageOpts string
		Field       string
		FieldOpts   string
		problems    testutils.Problems
	}{
		{"ValidCreateTime", messageOptsResource, `google.protobuf.Timestamp create_time`, fieldOptsOutputOnly, nil},
		{"InvalidCreateTime", messageOptsResource, `google.protobuf.Timestamp create_time`, ``, missingOutputOnly},
		{"ValidDeleteTime", messageOptsResource, `google.protobuf.Timestamp delete_time`, fieldOptsOutputOnly, nil},
		{"InvalidDeleteTime", messageOptsResource, `google.protobuf.Timestamp delete_time`, ``, missingOutputOnly},
		{"ValidUid", messageOptsResource, `string uid`, fieldOptsOutputOnly, nil},
		{"InvalidUid", messageOptsResource, `string uid`, ``, missingOutputOnly},
		{"ValidUpdateTime", messageOptsResource, `google.protobuf.Timestamp update_time`, fieldOptsOutputOnly, nil},
		{"InvalidUpdateTime", messageOptsResource, `google.protobuf.Timestamp update_time`, ``, missingOutputOnly},
		{"IrrelevantNotResource", ``, `string uid`, ``, nil},
		{"IrrelevantField", messageOptsResource, `string name`, ``, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";
				import "google/protobuf/timestamp.proto";
				message Book {
					{{.MessageOpts}}
					{{.Field}} = 1 {{.FieldOpts}};
				}
			`, test)
			field := file.Messages()[0].Fields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(fieldBehavior.Lint(file)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
