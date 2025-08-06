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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

// Tests our regexp normalizes strings to the expected path.
func TestVersionNormalization(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "foo.bar.baz",
			want: "",
		},
		{
			// This one's a bit iffy.  Should a version be allowed as the first segment?
			in:   "v1beta",
			want: "",
		},
		{
			in:   "foo.v3",
			want: "foo.v3",
		},
		{
			in:   "foo.v99alpha.bar",
			want: "foo.v99alpha",
		},
		{
			in:   "foo.v2.bar.v2",
			want: "foo.v2.bar.v2",
		},
	} {
		got := versionedPrefix.FindString(tc.in)
		if got != tc.want {
			t.Errorf("mismatch: in %q, got %q, want %q", tc.in, got, tc.want)
		}
	}
}

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
			description:   "unversioned subpackage",
			CallingPkg:    "somepackage",
			ReferencePkg:  "somepackage.sub",
			ReferencedMsg: "somepackage.sub.OtherMessage",
			problems:      testutils.Problems{{Message: "foreign type referenced"}},
		},
		{
			description:   "versioned subpackage",
			CallingPkg:    "somepackage.v6",
			ReferencePkg:  "somepackage.v6.sub",
			ReferencedMsg: "somepackage.v6.sub.OtherMessage",
			problems:      nil,
		},
		{
			description:   "versioned deep subpackaging",
			CallingPkg:    "somepackage.v1.abc",
			ReferencePkg:  "somepackage.v1.lol.xyz",
			ReferencedMsg: "somepackage.v1.lol.xyz.OtherMessage",
			problems:      nil,
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
			field := file.Messages().Get(0).Fields().Get(1)
			if diff := tc.problems.SetDescriptor(field).Diff(foreignTypeReference.Lint(file)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
