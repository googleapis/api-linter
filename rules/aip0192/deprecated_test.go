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

package aip0192

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

// These are split up since templating doesn't play nicely with inserting protobuf options.
func TestValidDescriptor(t *testing.T) {
	file := testutils.ParseProto3String(t, `
    // A library service.
    service Library {
      // Retrieves a book.
      rpc GetBook(GetBookRequest) returns (Book);
    }
    message GetBookRequest {}
    message Book {}
  `)

	serviceProblems := testutils.Problems{}
	if diff := serviceProblems.Diff(deprecatedComment.Lint(file)); diff != "" {
		t.Error(diff)
	}

	methodProblems := testutils.Problems{}
	if diff := methodProblems.Diff(deprecatedComment.Lint(file)); diff != "" {
		t.Error(diff)
	}
}

func TestDeprecatedMethod(t *testing.T) {
	tests := []struct {
		testName      string
		MethodComment string
		problems      testutils.Problems
	}{
		{"ValidMethodDeprecated", "// Deprecated: Don't use this.\n// Method comment.", nil},
		{"InvalidMethodDeprecated", "// Method comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
        service Library {

          {{.MethodComment}}
          rpc GetBook(GetBookRequest) returns (Book) {
            option deprecated = true;
          }
        }
        message GetBookRequest {}
        message Book {}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Services()[0].Methods()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeprecatedService(t *testing.T) {
	tests := []struct {
		testName       string
		ServiceComment string
		problems       testutils.Problems
	}{
		{"ValidServiceDeprecated", "// Deprecated: Don't use this.\n// Service comment.", nil},
		{"InvalidServiceDeprecated", "// Service comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
        {{.ServiceComment}}
        service Library {
          option deprecated = true;
          rpc GetBook(GetBookRequest) returns (Book);
        }
        message GetBookRequest {}
        message Book {}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Services()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeprecatedField(t *testing.T) {
	tests := []struct {
		testName     string
		FieldComment string
		problems     testutils.Problems
	}{
		{"ValidFieldDeprecated", "// Deprecated: Don't use this.\n// Field comment.", nil},
		{"InvalidFieldDeprecated", "// Field comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
        message GetBookRequest {
			{{.FieldComment}}
			string name = 1 [deprecated = true];
		}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Messages()[0].Fields()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeprecatedEnum(t *testing.T) {
	tests := []struct {
		testName    string
		EnumComment string
		problems    testutils.Problems
	}{
		{"ValidEnumDeprecated", "// Deprecated: Don't use this.\n// Enum comment.", nil},
		{"InvalidEnumDeprecated", "// Enum comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
		{{.EnumComment}}
		enum State {
			option deprecated = true;
			
			STATE_UNSPECIFIED = 0;
		}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Enums()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeprecatedEnumValue(t *testing.T) {
	tests := []struct {
		testName         string
		EnumValueComment string
		problems         testutils.Problems
	}{
		{"ValidEnumValueDeprecated", "// Deprecated: Don't use this.\n// EnumValue comment.", nil},
		{"InvalidEnumValueDeprecated", "// EnumValue comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
		enum State {
			{{.EnumValueComment}}
			STATE_UNSPECIFIED = 0 [deprecated = true];
		}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Enums()[0].Values()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeprecatedMessage(t *testing.T) {
	tests := []struct {
		testName       string
		MessageComment string
		problems       testutils.Problems
	}{
		{"ValidMessageDeprecated", "// Deprecated: Don't use this.\n// Message comment.", nil},
		{"InvalidMessageDeprecated", "// Message comment.", testutils.Problems{{Message: `Use "Deprecated: <reason>"`}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
		{{.MessageComment}}
        message GetBookRequest {
			option deprecated = true;

			string name = 1;
		}
      `, test)

			problems := deprecatedComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.Messages()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
