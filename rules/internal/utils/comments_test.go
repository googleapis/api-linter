// Copyright 2023 Google LLC
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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSeparateInternalComments(t *testing.T) {
	for _, tst := range []struct {
		name         string
		in           string
		wantExternal string
		wantInternal string
	}{
		{
			"external only",
			"Hello,\nWorld!",
			"Hello,\nWorld!",
			"",
		},
		{
			"internal only",
			"(-- Hello,\nInternal! --)",
			"",
			"Hello,\nInternal!",
		},
		{
			"mixed",
			"Hello,\nWorld!\n(-- We come\nin peace --)\nWhat planet is this?",
			"Hello,\nWorld!\nWhat planet is this?",
			"We come\nin peace",
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			got := SeparateInternalComments(tst.in)
			// Join them for ease of diffing.
			gotExternal := strings.Join(got.External, "\n")
			gotInternal := strings.Join(got.Internal, "\n")

			if diff := cmp.Diff(gotExternal, tst.wantExternal); diff != "" {
				t.Errorf("External: got(-),want(+):\n%s", diff)
			}
			if diff := cmp.Diff(gotInternal, tst.wantInternal); diff != "" {
				t.Errorf("Internal: got(-),want(+):\n%s", diff)
			}
		})
	}
}
