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

package utils

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestGetScalarTypeName(t *testing.T) {
	for _, ty := range []string{"int32", "int64", "string", "bytes", "float", "double"} {
		t.Run(ty, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message Book {
					{{.Type}} field = 1;
				}
			`, struct{ Type string }{ty})
			field := file.GetMessageTypes()[0].GetFields()[0]
			if got, want := GetScalarTypeName(field), ty; got != want {
				t.Errorf("Got %q, expected %q.", got, want)
			}
		})
	}
}

func TestGetMessageTypeName(t *testing.T) {
	tests := []struct {
		testName string
		Type     string
		want     string
	}{
		{"Message", "google.protobuf.Timestamp", "google.protobuf.Timestamp"},
		{"Scalar", "int32", ""},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/timestamp.proto";
				message Book {
					{{.Type}} field = 1;
				}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			if got := GetMessageTypeName(field); got != test.want {
				t.Errorf("Got %q, expected %q.", got, test.want)
			}
		})
	}
}
