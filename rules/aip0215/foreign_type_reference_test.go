// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0215

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestForeignTypeReference(t *testing.T) {

	for _, tc := range []struct {
		description   string
		CallingPkg    string
		ReferencePkg  string
		ReferencedMsg string
		problems      testutils.Problems
	}{
		{
			description:   "same pkg",
			CallingPkg:    "same",
			ReferencePkg:  "same",
			ReferencedMsg: "OtherMessage",
			problems:      nil,
		},
		{
			description:   "refers to google.api",
			CallingPkg:    "same",
			ReferencePkg:  "google.api",
			ReferencedMsg: "google.api.OtherMessage",
			problems:      nil,
		},
		{
			description:   "refers to google.protobuf",
			CallingPkg:    "same",
			ReferencePkg:  "google.protobuf",
			ReferencedMsg: "google.protobuf.OtherMessage",
			problems:      nil,
		},
		{
			description:   "refers to foreign pkg",
			CallingPkg:    "same",
			ReferencePkg:  "other",
			ReferencedMsg: "other.OtherMessage",
			problems:      testutils.Problems{{Message: "foreign type referenced"}},
		},
		{
			description:   "refers to subpkg",
			CallingPkg:    "somepackage",
			ReferencePkg:  "somepackage.sub",
			ReferencedMsg: "somepackage.sub.OtherMessage",
			problems:      testutils.Problems{{Message: "foreign type referenced"}},
		},
		{
			description:   "refers to component package",
			CallingPkg:    "somepackage",
			ReferencePkg:  "otherpackage.type",
			ReferencedMsg: "otherpackage.type.OtherMessage",
			problems:      nil,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			files := testutils.ParseProto3Tmpls(t, map[string]string{
				"calling.proto": `
					package {{.CallingPkg}};
					import "ref.proto";
					message Caller {
						string foo = 1;
						{{.ReferencedMsg}} bar = 2;
					}
				`,
				"ref.proto": `
					package {{.ReferencePkg}};
					message OtherMessage {
						int32 baz = 1;
					}
				`,
			}, tc)
			file := files["calling.proto"]
			field := file.GetMessageTypes()[0].GetFields()[1]
			if diff := tc.problems.SetDescriptor(field).Diff(foreignTypeReference.Lint(file)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
