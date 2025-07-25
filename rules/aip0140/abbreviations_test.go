// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

func TestAbbreviations(t *testing.T) {
	ruleGroups := []struct {
		name     string
		tmpl     string
		caseFunc func(string) string
		descFunc func(protoreflect.FileDescriptor) protoreflect.Descriptor
	}{
		{
			name:     "Field",
			tmpl:     `message Book { string {{.Name}} = 1; }`,
			caseFunc: strcase.SnakeCase,
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Messages().Get(0).Fields().Get(0)
			},
		},
		{
			name:     "Message",
			tmpl:     `message {{.Name}} {}`,
			caseFunc: strcase.UpperCamelCase,
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Messages().Get(0)
			},
		},
		{
			name:     "Service",
			tmpl:     `service {{.Name}} {}`,
			caseFunc: strcase.UpperCamelCase,
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Services().Get(0)
			},
		},
		{
			name: "Method",
			tmpl: `
				service Library {
					rpc {{.Name}}(Request) returns (Response);
				}

				message Request {}
				message Response {}
			`,
			caseFunc: strcase.UpperCamelCase,
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Services().Get(0).Methods().Get(0)
			},
		},
		{
			name:     "Enum",
			tmpl:     `enum {{.Name}} { UNSPECIFIED = 0; }`,
			caseFunc: strcase.UpperCamelCase,
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Enums().Get(0)
			},
		},
		{
			name: "EnumValue",
			tmpl: `enum Thing { {{.Name}} = 0; }`,
			caseFunc: func(s string) string {
				return strings.ToUpper(strcase.SnakeCase(s))
			},
			descFunc: func(f protoreflect.FileDescriptor) protoreflect.Descriptor {
				return f.Enums().Get(0).Values().Get(0)
			},
		},
	}
	for _, group := range ruleGroups {
		t.Run(group.name, func(t *testing.T) {
			for _, test := range buildTests(group.caseFunc) {
				t.Run(test.Name, func(t *testing.T) {
					// Build the file descriptor, and retrieve the descriptor we
					// expect any problems to be attached to.
					file := testutils.ParseProto3Tmpl(t, group.tmpl, test)
					d := group.descFunc(file)

					// Establish that we get the problems we expect.
					problems := abbreviations.Lint(file)
					if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
						t.Error(diff)
					}
				})
			}
		})
	}
}

type abbvTest struct {
	Name     string
	problems testutils.Problems
}

func buildTests(caseFunc func(string) string) []abbvTest {
	return []abbvTest{
		{caseFunc("book_configuration"), testutils.Problems{{Suggestion: caseFunc("book_config")}}},
		{caseFunc("book_config"), testutils.Problems{}},
		{caseFunc("book_identifier"), testutils.Problems{{Suggestion: caseFunc("book_id")}}},
		{caseFunc("book_id"), testutils.Problems{}},
		{caseFunc("book_information"), testutils.Problems{{Suggestion: caseFunc("book_info")}}},
		{caseFunc("book_info"), testutils.Problems{}},
		{caseFunc("book_specification"), testutils.Problems{{Suggestion: caseFunc("book_spec")}}},
		{caseFunc("book_spec"), testutils.Problems{}},
		{caseFunc("book_statistics"), testutils.Problems{{Suggestion: caseFunc("book_stats")}}},
		{caseFunc("book_stats"), testutils.Problems{}},
		{caseFunc("informational_book"), testutils.Problems{}},
	}
}