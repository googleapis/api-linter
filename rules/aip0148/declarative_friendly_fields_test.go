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
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestDeclarativeFriendlyFields(t *testing.T) {
	for _, test := range []struct {
		name    string
		skipped stringset.Set
	}{
		{"Valid", stringset.New()},
		{"Name", stringset.New("name")},
		{"UID", stringset.New("uid")},
		{"DisplayName", stringset.New("display_name")},
		{"CreateTime", stringset.New("create_time")},
		{"UpdateTime", stringset.New("update_time")},
		{"DeleteTime", stringset.New("delete_time")},
		{"AllTimes", stringset.New("create_time", "update_time", "delete_time")},
		{"Randos", stringset.New("uid", "display_name")},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Set up the string with the fields we will include.
			fields := ""
			cursor := 1
			for fieldName, fieldType := range reqFields {
				if !test.skipped.Contains(fieldName) {
					fields += fmt.Sprintf("  %s %s = %d;\n", fieldType, fieldName, cursor)
					cursor++
				}
			}

			// Create the potential problem object for the missing fields.
			var problems testutils.Problems
			if test.skipped.Len() == 1 {
				f := test.skipped.Unordered()[0]
				problems = testutils.Problems{{
					Message: fmt.Sprintf("must include the `%s %s` field", reqFields[f], f),
				}}
			} else if test.skipped.Len() > 1 {
				missingFields := stringset.New()
				for _, f := range test.skipped.Unordered() {
					missingFields.Add(fmt.Sprintf("%s %s", reqFields[f], f))
				}
				msg := ""
				for _, f := range missingFields.Elements() {
					msg += fmt.Sprintf("  - `%s`\n", f)
				}
				problems = testutils.Problems{{Message: strings.TrimSuffix(msg, "\n")}}
			}

			// Test against declarative-friendly and standard styles.
			for _, subtest := range []struct {
				name     string
				style    string
				problems testutils.Problems
			}{
				{"DeclFriendly", "style: DECLARATIVE_FRIENDLY", problems},
				{"NotDeclFriendly", "", nil},
			} {
				t.Run(subtest.name, func(t *testing.T) {
					f := testutils.ParseProto3Tmpl(t, `
						import "google/api/resource.proto";
						import "google/protobuf/timestamp.proto";
						message Book {
							option (google.api.resource) = {
								type: "library.googleapis.com/Book"
								pattern: "publishers/{publisher}/books/{book}"
								{{.Style}}
							};
							{{.Fields}}
						}
					`, struct {
						Fields string
						Style  string
					}{Fields: fields, Style: subtest.style})
					m := f.Messages().Get(0)
					got := declarativeFriendlyRequired.Lint(f)
					if diff := subtest.problems.SetDescriptor(m).Diff(got); diff != "" {
						t.Error(diff)
					}
				})
			}
		})
	}
}

func TestDeclarativeFriendlyFieldsSingleton(t *testing.T) {
	for _, test := range []struct {
		name   string
		Fields string
		want   testutils.Problems
	}{
		{
			"InvalidNoCreateTime", `string name = 1; string display_name = 2; google.protobuf.Timestamp update_time = 3;`,
			testutils.Problems{{Message: "create_time"}},
		},
		{
			"ValidNoDeleteTimeNoUid", `string name = 1; string display_name = 2; ` +
				`google.protobuf.Timestamp create_time = 3; google.protobuf.Timestamp update_time = 4;`,
			nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/timestamp.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Settings"
						pattern: "publishers/{publisher}/settings"
						style: DECLARATIVE_FRIENDLY
					};
					{{.Fields}}
				}
			`, test)
			m := f.Messages().Get(0)
			got := declarativeFriendlyRequired.Lint(f)
			if diff := test.want.SetDescriptor(m).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
