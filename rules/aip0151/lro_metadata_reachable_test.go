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

package aip0151

import (
	"strings"
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestLROMetadataReachable(t *testing.T) {
	t.Run("SameFile", func(t *testing.T) {
		tests := []struct {
			testName     string
			Package      string
			MetadataType string
			problems     testutils.Problems
		}{
			{"Valid", "", "WriteBookMetadata", testutils.Problems{}},
			{"ValidPkg", "package test;", "WriteBookMetadata", testutils.Problems{}},
			{"InvalidTypo", "", "WriteBookMeta", testutils.Problems{{Message: "WriteBookMeta"}}},
			{"InvalidTypoPkg", "package test;", "WriteBookMeta", testutils.Problems{{Message: "test.WriteBookMeta"}}},
			{"ValidExternal", "", "google.protobuf.Empty", testutils.Problems{}},
			{"ValidExternalPkg", "package test;", "google.protobuf.Empty", testutils.Problems{}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					{{.Package}}
					import "google/longrunning/operations.proto";
					service Library {
						rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
							option (google.longrunning.operation_info) = {
								response_type: "WriteBookResponse"
								metadata_type: "{{.MetadataType}}"
							};
						}
					}
					message WriteBookRequest {}
					message WriteBookResponse {}
					message WriteBookMetadata {}
				`, test)
				problems := lroMetadataReachable.Lint(f)
				m := f.GetServices()[0].GetMethods()[0]
				if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
					t.Errorf(diff)
				}

			})
		}
	})
	t.Run("Imported", func(t *testing.T) {
		for _, test := range []struct {
			name     string
			message  string
			problems testutils.Problems
		}{
			{"Present", "message WriteBookMetadata {}", testutils.Problems{}},
			{"Absent", "", testutils.Problems{{Message: "WriteBookMetadata"}}},
		} {
			t.Run(test.name, func(t *testing.T) {
				files := testutils.ParseProtoStrings(t, map[string]string{
					"imported.proto": strings.ReplaceAll(`
						syntax = "proto3";
						message WriteBookRequest {}
						message WriteBookResponse {}
						---
					`, "---", test.message),
					"test.proto": `
						syntax = "proto3";
						import "google/longrunning/operations.proto";
						import "imported.proto";
						service Library {
							rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
								option (google.longrunning.operation_info) = {
									response_type: "WriteBookResponse"
									metadata_type: "WriteBookMetadata"
								};
							}
						}
					`,
				})
				problems := lroMetadataReachable.Lint(files["test.proto"])
				method := files["test.proto"].GetServices()[0].GetMethods()[0]
				if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})
}
