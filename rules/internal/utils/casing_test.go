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

func TestIsLowerCamelCase(t *testing.T) {
	for _, test := range []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Valid",
			s:    "bookShelf123",
			want: true,
		},
		{
			name: "InvalidUpperCamelCase",
			s:    "BookShelf",
			want: false,
		},
		{
			name: "InvalidDash",
			s:    "book-Shelf",
			want: false,
		},
		{
			name: "InvalidUnderscore",
			s:    "book_Shelf",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := IsLowerCamelCase(test.s)
			if got != test.want {
				t.Errorf(
					"IsLowerCamelCase(%q) = %v, want %v",
					test.s, got, test.want,
				)
			}
		})
	}
}

func TestIsUpperCamelCase(t *testing.T) {
	for _, test := range []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Valid",
			s:    "BookShelf123",
			want: true,
		},
		{
			name: "InvalidLowerCamelCase",
			s:    "bookShelf",
			want: false,
		},
		{
			name: "InvalidDash",
			s:    "Book-Shelf",
			want: false,
		},
		{
			name: "InvalidUnderscore",
			s:    "Book_Shelf",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := IsUpperCamelCase(test.s)
			if got != test.want {
				t.Errorf(
					"IsLowerCamelCase(%q) = %v, want %v",
					test.s, got, test.want,
				)
			}
		})
	}
}
