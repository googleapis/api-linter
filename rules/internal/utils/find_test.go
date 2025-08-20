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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

// func TestFindMessage(t *testing.T) {
// 	files := testutils.CompileStrings(t, map[string]string{
// 		"a.proto": `
// 			syntax = "proto3";
// 			package test;
// 			message Book {}
// 		`,
// 		"b.proto": `
// 			syntax = "proto3";
// 			package other;
// 			message Scroll {}
// 		`,
// 		"c.proto": `
// 			syntax = "proto3";
// 			package test;
// 			import "a.proto";
// 			import "b.proto";
// 		`,
// 	})
// 	if book := FindMessage(files[2], "Book"); book == nil {
// 		t.Errorf("Got nil, expected Book message.")
// 	}
// 	if scroll := FindMessage(files[2], "Scroll"); scroll != nil {
// 		t.Errorf("Got Scroll message, expected nil.")
// 	}
// 	if book := FindMessage(files[2], "test.Book"); book == nil {
// 		t.Errorf("Got nil, expected Book message from qualified name.")
// 	}
// 	if scroll := FindMessage(files[2], "other.Scroll"); scroll == nil {
// 		t.Errorf("Got nil message, expected Scroll message from qualified name.")
// 	}
// }

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
	msg := file.Messages().Get(0)

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
			} else if got := f.Name(); string(got) != want {
				t.Errorf("Got %q, expected %q", got, want)
			}
		})
	}

}
