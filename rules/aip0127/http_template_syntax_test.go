// Copyright 2022 Google LLC
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

package aip0127

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpTemplateSyntax(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		valid    bool
	}{
		// NOTE: These examples are contrived to test the enforcment of the
		// template syntax. Many of these examples either fail or do not make
		// sense in the context of other AIP rules.

		// Valid cases
		{"SingleLiteral", "/v1", true},
		{"VersionedTemplate", "/{$api_version}/books", true},
		{"TwoLiterals", "/v1/books", true},
		{"ThreeLiterals", "/v1/books/shelves", true},
		{"SingleLiteralWithVerb", "/v1:verb", true},
		{"MultipleLiteralsWithVerb", "/v1/books:verb", true},
		{"SingleWildcard", "/v1/*", true},
		{"DoubleWildcard", "/v1/**", true},
		{"SingleWildcardWithVerb", "/v1/*:verb", true},
		{"DoubleWildcardWithVerb", "/v1/**:verb", true},
		{"SingleWildcardFollowedByLiteral", "/v1/*/books", true},
		{"DoubleWildcardFollowedByLiteral", "/v1/**/books", true},
		{"LiteralFollowedBySingleWildcard", "/v1/books/*", true},
		{"LiteralFollowedByDoubleWildcard", "/v1/books/**", true},
		{"VariableWithFieldpath", "/v1/{field}", true},
		{"VariableWithNestedFieldpath", "/v1/{field.subfield}", true},
		{"VariableWithUltraNestedFieldpath", "/v1/{field.subfield.subsubfield}", true},
		{"VariableWithLiteralTemplate", "/v1/{field=books}", true},
		{"VariableWithSingleWildcardTemplate", "/v1/{field=*}", true},
		{"VariableWithDoubleWildcardTemplate", "/v1/{field=**}", true},
		{"VariableWithSingleWildcardFollowedByLiteral", "/v1/{field=*/books}", true},
		{"VariableWithDoubleWildcardFollowedByLiteral", "/v1/{field=**/books}", true},
		{"VariableWithLiteralFollowedBySingleWildcard", "/v1/{field=books/*}", true},
		{"VariableWithLiteralFollowedByDoubleWildcard", "/v1/{field=books/**}", true},
		{"VariableFollowedByLiteral", "/v1/{field}/books", true},
		{"VariableFollowedByVariable", "/v1/{field}/{otherField}", true},
		{"VariableWithTemplateFollowedByLiteral", "/v1/{field=books/*}/shelves", true},
		{"VariableFollowedByVariableWithTemplate", "/v1/{field}/{otherField=books/*}", true},
		{"VariableWithTemplateFollowedByVariableWithTemplate", "/v1/{field=books/*}/{otherField=shelves/*}", true},

		// Invalid cases
		{"LiteralWithoutLeadingSlash", "v1", false},
		{"LiteralFollowedBySlash", "/v1/", false},
		{"WrongVerbDelimiter", "/v1-verb", false},
		{"MultipleVerbs", "/v1:verb:verb", false},
		{"VerbFollowedBySlash", "/v1:verb/", false},
		{"MultipleLiteralsWithWrongDelimiter", "/v1|books", false},
		{"SingleWildcardFollowedBySlash", "/v1/*/", false},
		{"DoubleWildcardFollowedBySlash", "/v1/**/", false},
		{"TripleWildcard", "/v1/***", false},
		{"WrongVariableMarker", "/v1/[field]", false},
		{"WrongVariableSubfieldOperator", "/v1/[field->subfield]", false},
		{"VariableTemplateContainsVariable", "/v1/{field={otherField=*}}", false},
		{"WrongVariableTemplateAssignmentOperator", "/v1/{field≈books}", false},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc FooMethod(FooMethodRequest) returns (FooMethodResponse) {
						option (google.api.http) = {
							get: "{{.URI}}"
						};
					}
				}
				message FooMethodRequest {}
				message FooMethodResponse {}
			`, test)

			problems := httpTemplateSyntax.Lint(file)

			if test.valid && len(problems) > 0 {
				t.Fatalf("Expected valid HTTP path template syntax but got invalid")
			}

			if !test.valid && len(problems) == 0 {
				t.Fatalf("Expected invalid HTTP path template syntax but got valid")
			}
		})
	}
}
