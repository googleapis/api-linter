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

package aip0192

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestAbsoluteLinks(t *testing.T) {
	for _, test := range []struct {
		name     string
		URI      string
		problems testutils.Problems
	}{
		{"Valid", "https://google.com/", testutils.Problems{}},
		{"NoProtocol", "google.com", testutils.Problems{{Message: "https://"}}},
		{"Relative", "/foo/bar/baz.html", testutils.Problems{{Message: "https://"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			  // This is [a link]({{.URI}}).
				message Foo {}
			`, test)
			m := f.GetMessageTypes()[0]
			problems := absoluteLinks.Lint(f)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
