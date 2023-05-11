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

import "testing"

func TestToLowerCamelCase(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OneWord",
			input: "Foo",
			want:  "foo",
		},
		{
			name:  "OneWordNoop",
			input: "foo",
			want:  "foo",
		},
		{
			name:  "TwoWords",
			input: "bookShelf",
			want:  "bookShelf",
		},
		{
			name:  "WithDash",
			input: "book-shelf",
			want:  "bookShelf",
		},
		{
			name:  "WithNumbers",
			input: "universe42love",
			want:  "universe42love",
		},
		{
			name:  "WithUnderscore",
			input: "book_shelf",
			want:  "bookShelf",
		},
		{
			name:  "WithUnderscore",
			input: "book_shelf",
			want:  "bookShelf",
		},
		{
			name:  "WithSpaces",
			input: "book shelf",
			want:  "bookShelf",
		},
		{
			name:  "WithPeriods",
			input: "book.shelf",
			want:  "bookShelf",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := ToLowerCamelCase(test.input)
			if got != test.want {
				t.Errorf("ToLowerCamelCase(%q) = %q, got %q", test.input, test.want, got)
			}
		})
	}
}

func TestToUpperCamelCase(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OneWord",
			input: "foo",
			want:  "Foo",
		},
		{
			name:  "OneWordNoop",
			input: "Foo",
			want:  "Foo",
		},
		{
			name:  "TwoWords",
			input: "bookShelf",
			want:  "BookShelf",
		},
		{
			name:  "WithDash",
			input: "book-shelf",
			want:  "BookShelf",
		},
		{
			name:  "WithNumbers",
			input: "universe42love",
			want:  "Universe42love",
		},
		{
			name:  "WithUnderscore",
			input: "Book_shelf",
			want:  "BookShelf",
		},
		{
			name:  "WithUnderscore",
			input: "Book_shelf",
			want:  "BookShelf",
		},
		{
			name:  "WithSpaces",
			input: "Book shelf",
			want:  "BookShelf",
		},
		{
			name:  "WithPeriods",
			input: "book.shelf",
			want:  "BookShelf",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := ToUpperCamelCase(test.input)
			if got != test.want {
				t.Errorf("ToLowerCamelCase(%q) = %q, got %q", test.input, test.want, got)
			}
		})
	}
}
