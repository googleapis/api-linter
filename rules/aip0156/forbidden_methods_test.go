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

package aip0156

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestForbiddenMethods(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		Suffix     string
		problems   testutils.Problems
	}{
		{"ValidGet", "GetPublisherSettings", "/settings", testutils.Problems{}},
		{"ValidUpdate", "UpdatePublisherSettings", "/settings", testutils.Problems{}},
		{"InvalidCreate", "CreatePublisherSettings", "/settings", testutils.Problems{{Message: "Create"}}},
		{"InvalidList", "ListPublisherSettings", "/settings", testutils.Problems{{Message: "List"}}},
		{"InvalidDelete", "DeletePublisherSettings", "/settings", testutils.Problems{{Message: "Delete"}}},
		{"Irrelevant", "CreatePublisher", "", testutils.Problems{}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/api/annotations.proto";
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns (Settings) {
					option (google.api.http) = {
						post: "/v1/{name=publishers/*{{.Suffix}}}"
						body: "settings"
					};
				}
			}
			message {{.MethodName}}Request {}
			message Settings{}
		`, test)
		method := f.GetServices()[0].GetMethods()[0]
		problems := forbiddenMethods.Lint(method)
		if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
			t.Errorf(diff)
		}
	}
}
