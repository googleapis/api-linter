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

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestFindMessage(t *testing.T) {
	files := testutils.ParseProtoStrings(t, map[string]string{
		"a.proto": `
			package test;
			message Book {}
		`,
		"b.proto": `
			package other;
			message Scroll {}
		`,
		"c.proto": `
			package test;
			import "a.proto";
			import "b.proto";
		`,
	})
	if book := FindMessage(files["c.proto"], "Book"); book == nil {
		t.Errorf("Got nil, expected Book message.")
	}
	if scroll := FindMessage(files["c.proto"], "Scroll"); scroll != nil {
		t.Errorf("Got Sctoll message, expected nil.")
	}
}
