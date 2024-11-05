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
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestAnnotationExistsValid(t *testing.T) {
	f := testutils.ParseProto3String(t, `
		import "google/longrunning/operations.proto";
		service Library {
			rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "WriteBookResponse"
					metadata_type: "WriteBookMetadata"
				};
			}
		}
		message WriteBookRequest {}
	`)
	if diff := (testutils.Problems{}).Diff(lroAnnotationExists.Lint(f)); diff != "" {
		t.Error(diff)
	}
}

func TestAnnotationExistsInvalid(t *testing.T) {
	f := testutils.ParseProto3String(t, `
		import "google/longrunning/operations.proto";
		service Library {
			rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation);
		}
		message WriteBookRequest {}
	`)
	want := testutils.Problems{{
		Descriptor: f.GetServices()[0].GetMethods()[0],
		Message:    "operation_info annotation",
	}}
	if diff := want.Diff(lroAnnotationExists.Lint(f)); diff != "" {
		t.Error(diff)
	}
}
