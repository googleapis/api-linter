// Copyright 2020 Google LLC
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

func TestIsCommonProto(t *testing.T) {
	for _, test := range []struct {
		Package string
		want    bool
	}{
		{"google.api", true},
		{"google.longrunning", true},
		{"google.protobuf", true},
		{"google.rpc", true},
		{"google.api.experimental", true},
		{"google.cloud.speech.v1", false},
	} {
		t.Run(test.Package, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				package {{.Package}};
			`, test)
			if got := IsCommonProto(f); got != test.want {
				t.Errorf("Got %v, expected %v", got, test.want)
			}
		})
	}
}
