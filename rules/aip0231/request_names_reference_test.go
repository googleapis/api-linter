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

package aip0231

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestRequestNamesReference(t *testing.T) {
	for _, test := range []struct {
		name        string
		MessageName string
		FieldName   string
		FieldOpts   string
		problems    testutils.Problems
	}{
		{"Valid", "BatchGetBooks", "names", ` [(google.api.resource_reference) = { type: "foo" }]`, testutils.Problems{}},
		{"Missing", "BatchGetBooks", "names", "", testutils.Problems{{Message: "google.api.resource_reference"}}},
		{"IrrelevantMessage", "GetBooks", "names", "", testutils.Problems{}},
		{"IrrelevantField", "BatchGetBooks", "ids", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message {{.MessageName}}Request {
					repeated string {{.FieldName}} = 1{{.FieldOpts}};
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(requestNamesReference.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
