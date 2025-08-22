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

package aip0191

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestPhpNamespace(t *testing.T) {
	for _, test := range []struct {
		name         string
		PhpNamespace string
		problems     testutils.Problems
	}{
		{"Valid", `Google\\Example\\V1`, testutils.Problems{}},
		{"ValidBeta", `Google\\Example\\V1beta1`, testutils.Problems{}},
		{"ValidServiceNamespaceCase", `Google\\FooBar\\V1`, testutils.Problems{}},
		{"InvalidBadChars", "Google:Example:V1", testutils.Problems{{Message: "Invalid characters"}}},
		{"Invalid", `google\\example\\v1`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1`),
		}}},
		{"InvalidVersion", `Google\\Example\\v1`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1`),
		}}},
		{"InvalidBeta", `Google\\Example\\V1Beta1`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1beta1`),
		}}},
		{"InvalidBetaChannel", `Google\\Example\\V1Beta`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1beta`),
		}}},
		{"InvalidTemplate", `Google\\Example\\V1Main`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1main`),
		}}},
		{"InvalidPointRelease", `Google\\Example\\V1P1Beta1`, testutils.Problems{{
			Suggestion: fmt.Sprintf("option php_namespace = %s;", `Google\\Example\\V1p1beta1`),
		}}},
		{"ValidServiceNamespaceCase", `Google\\Foobar\\V1`, testutils.Problems{{Message: "Case"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				package google.example.v1;

				option php_namespace = "{{.PhpNamespace}}";

				service FooBar {}
			`, test)
			if diff := test.problems.SetDescriptor(f).Diff(phpNamespace.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
