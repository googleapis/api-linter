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
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
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

func TestFindFieldDotNotation(t *testing.T) {
	file := testutils.ParseProto3String(t, `
		package test;
		
		message CreateBookRequest {
			string parent = 1;

			Book book = 2;
		}

		message Book {
			string name = 1;

			message PublishingInfo {
				string publisher = 1;
				int32 edition = 2;
			}

			PublishingInfo publishing_info = 2;
		}
	`)
	msg := file.GetMessageTypes()[0]

	for _, tst := range []struct {
		name, path string
	}{
		{"top_level", "parent"},
		{"nested", "book.name"},
		{"double_nested", "book.publishing_info.publisher"},
	} {
		t.Run(tst.name, func(t *testing.T) {
			split := strings.Split(tst.path, ".")
			want := split[len(split)-1]

			f := FindFieldDotNotation(msg, tst.path)

			if f == nil {
				t.Errorf("Got nil, expected %q field", want)
			} else if got := f.GetName(); got != want {
				t.Errorf("Got %q, expected %q", got, want)
			}
		})
	}

}
