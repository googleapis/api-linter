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

package aip0136

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestURISuffix(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		URI        string
		problems   testutils.Problems
	}{
		{"ValidTwoWordNoParent", "SearchAudioBooks", "/v1/audioBooks:search", nil},
		{"ValidVerb", "ArchiveBook", "/v1/{resource_name=publishers/*/books}:archive", testutils.Problems{}},
		{"ValidVerbNested", "ArchiveBook", "/v1/{book.resource_name=publishers/*/books}:archive", testutils.Problems{}},
		{"ValidVerbParent", "ImportBooks", "/v1/{parent=publishers/*}/books:import", testutils.Problems{}},
		{"ValidVerbParentBatchGet", "BatchGetBooks", "/v1/{parent=publishers/*}/books:batchGet", testutils.Problems{}},
		{"InvalidVerb", "ArchiveBook", "/v1/{resource_name=publishers/*/books}:archiveBook", testutils.Problems{{Message: ":archive"}}},
		{"ValidVerbNounNoVars", "TranslateText", "/v3:translateText", testutils.Problems{}},
		{"ValidVerbNounNoName", "TranslateText", "/v3/{location=projects/*/locations/*}:translateText", testutils.Problems{}},
		{"InvalidVerbNoun", "TranslateText", "/v3:translate", testutils.Problems{{Message: ":translateText"}}},
		{"ValidOneWord", "Translate", "/v3:translate", testutils.Problems{}},
		{"ValidStdMethod", "GetBook", "/v1/{resource_name=publishers/*/books/*}", testutils.Problems{}},
		{"ValidTwoWordNoun", "WriteAudioBook", "/v1/{resource_name=publishers/*/audioBooks/*}:write", testutils.Problems{}},
		{"ValidListRevisions", "ListBookRevisions", "/v1/{resource_name=publishers/*/books/*}:listRevisions", testutils.Problems{}},
		{"ValidTagRevision", "TagBookRevision", "/v1/{resource_name=publishers/*/books/*}:tagRevision", testutils.Problems{}},
		{"ValidDeleteRevision", "DeleteBookRevision", "/v1/{resource_name=publishers/*/books/*}:deleteRevision", testutils.Problems{}},
		{"ValidCollection", "SortBooks", "/v1/{publisher=publishers/*}/books:sort", testutils.Problems{}},
		{"ValidNoParent", "SearchBooks", "/v1/books:search", testutils.Problems{}},
		{"InvalidListRevisions", "ListBookRevisions", "/v1/{resource_name=publishers/*/books/*}:list", testutils.Problems{{Message: ":listRevisions"}}},
		{"InvalidTagRevision", "TagBookRevision", "/v1/{resource_name=publishers/*/books/*}:tag", testutils.Problems{{Message: ":tagRevision"}}},
		{"InvalidDeleteRevision", "DeleteBookRevision", "/v1/{resource_name=publishers/*/books/*}:delete", testutils.Problems{{Message: ":deleteRevision"}}},
		{"IgnoredFailsVariables", "AddPages", "/v1/{resource_name=publishers/*/books/*}:addPages", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							post: "{{.URI}}"
							body: "*"
						};
					}
				}
				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			problems := uriSuffix.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
