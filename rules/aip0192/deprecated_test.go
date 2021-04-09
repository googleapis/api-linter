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
func TestDeprecatedMethod(t *testing.T) {
	tests := []struct {
		testName      string
		MethodComment string
		problems      testutils.Problems
	}{
		{"ValidMethodDeprecated", "// Deprecated: Don't use this.\n// Method comment.", nil},
		{"InvalidMethodDeprecated", "// Method comment.", testutils.Problems{{Message: "Deprecated methods"}}},
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

			problems := deprecatedMethodComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.GetServices()[0].GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
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
		{"InvalidServiceDeprecated", "// Service comment.", testutils.Problems{{Message: "Deprecated services"}}},
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

			problems := deprecatedServiceComment.Lint(file)
			if diff := test.problems.SetDescriptor(file.GetServices()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
