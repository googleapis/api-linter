// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestListResponseResourceName(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want string
	}{
		{"ValidBooks", "ListBooks", "books"},
		{"ValidCamelCase", "ListBookShelves", "book_shelves"},
		{"InvalidListRevisions", "ListBookRevisions", ""},
		{"InvalidNotList", "WriteBook", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Response {}
			`, test)
			m := file.Messages().Get(0)
			got := ListResponseResourceName(m)
			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("got(-),want(+):\n%s", diff)
			}

		})
	}
}

func TestIsListResponseMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid list response",
			RPC:  "ListBooks",
			want: true,
		},
		{
			name: "not list response",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike response",
			RPC:  "Listen",
			want: false,
		},
		{
			name: "multiword lookalike response",
			RPC:  "ListenForever",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Response {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsListResponseMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsListRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid list request",
			RPC:  "ListBooks",
			want: true,
		},
		{
			name: "not list request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Listen",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "ListenForever",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Request {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsListRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsListRevisionsResponseMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid list revisions response",
			RPC:  "ListBook",
			want: true,
		},
		{
			name: "not list revisions response",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike response",
			RPC:  "Listen",
			want: false,
		},
		{
			name: "multiword lookalike response",
			RPC:  "ListenForever",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}RevisionsResponse {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsListRevisionsResponseMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsListRevisionsRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid list revisions request",
			RPC:  "ListBook",
			want: true,
		},
		{
			name: "not list revisions request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Listen",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "ListenForever",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}RevisionsRequest {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsListRevisionsRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsGetRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid get request",
			RPC:  "GetBook",
			want: true,
		},
		{
			name: "not get request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Getup",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "GetupFinder",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Request {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsGetRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsCreateRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid create request",
			RPC:  "CreateBook",
			want: true,
		},
		{
			name: "not create request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Created",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "CreatedBook",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Request {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsCreateRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsUpdateRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid update request",
			RPC:  "UpdateBook",
			want: true,
		},
		{
			name: "not update request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Updater",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "UpdaterAndCreator",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Request {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsUpdateRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsDeleteRequestMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPC  string
		want bool
	}{
		{
			name: "valid delete request",
			RPC:  "DeleteBook",
			want: true,
		},
		{
			name: "not delete request",
			RPC:  "ArchiveBook",
			want: false,
		},
		{
			name: "lookalike request",
			RPC:  "Deleter",
			want: false,
		},
		{
			name: "multiword lookalike request",
			RPC:  "DeleterAndPurge",
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.Compile(t, `
				message {{.RPC}}Request {}
			`, test)
			m := file.Messages().Get(0)
			if got, want := IsDeleteRequestMessage(m), test.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
