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

package aip0140

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestBase64(t *testing.T) {
	tmpl := `
		message Foo {
			// {{.Comment}}
		  {{.ProtoType}} bar = 1;
		}
	`
	errMsg := "Field \"bar\" mentions base64 encoding in comments, so it should probably be type `bytes`, not `string`."
	tests := []struct {
		testName  string
		Comment   string
		ProtoType string
		problems  testutils.Problems
	}{
		{"ValidNoBase64", "Blah blah blah.", "string", testutils.Problems{}},
		{"ValidBytes", "base64 encoded", "bytes", testutils.Problems{}},
		{"ValidBytesHyphen", "base-64 encoded", "bytes", testutils.Problems{}},
		{"Invalid", "base64 encoded", "string", testutils.Problems{{Message: errMsg}}},
		{"InvalidHyphen", "base-64 encoded", "string", testutils.Problems{{Message: errMsg}}},
		{"InvalidCaps", "Base64 encoded", "string", testutils.Problems{{Message: errMsg}}},
		{"InvalidCapsHyphen", "Base-64 encoded", "string", testutils.Problems{{Message: errMsg}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, tmpl, test)
			field := file.Messages().Get(0).Fields().Get(0)
			problems := base64.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
