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
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestGetResourceSingular(t *testing.T) {
	for _, test := range []struct {
		name     string
		resource *annotations.ResourceDescriptor
		want     string
	}{
		{
			name: "SingularSpecified",
			resource: &annotations.ResourceDescriptor{
				Singular: "bookShelf",
			},
			want: "bookShelf",
		},
		{
			name: "SingularAndTypeSpecified",
			resource: &annotations.ResourceDescriptor{
				Singular: "bookShelf",
				// NOTE: this is not a correct resource annotation.
				// it must match singular.
				Type: "library.googleapis.com/book",
			},
			want: "bookShelf",
		},
		{
			name: "TypeSpecified",
			resource: &annotations.ResourceDescriptor{
				Type: "library.googleapis.com/bookShelf",
			},
			want: "bookShelf",
		},
		{
			name:     "NothingSpecified",
			resource: &annotations.ResourceDescriptor{},
			want:     "",
		},
		{
			name:     "Nil",
			resource: nil,
			want:     "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := GetResourceSingular(test.resource)
			if got != test.want {
				t.Errorf("GetResourceSingular: expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestGetResourcePlural(t *testing.T) {
	for _, test := range []struct {
		name     string
		resource *annotations.ResourceDescriptor
		want     string
	}{
		{
			name: "PluralSpecified",
			resource: &annotations.ResourceDescriptor{
				Plural: "bookShelves",
			},
			want: "bookShelves",
		},
		{
			name:     "NothingSpecified",
			resource: &annotations.ResourceDescriptor{},
			want:     "",
		},
		{
			name:     "Nil",
			resource: nil,
			want:     "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := GetResourcePlural(test.resource)
			if got != test.want {
				t.Errorf("GetResourcePlural: expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestIsResourceRevision(t *testing.T) {
	for _, test := range []struct {
		name, Message, Resource string
		want                    bool
	}{
		{
			name:     "valid_revision",
			Message:  "BookRevision",
			Resource: `option (google.api.resource) = {type: "library.googleapis.com/BookRevision"};`,
			want:     true,
		},
		{
			name:    "not_revision_no_resource",
			Message: "BookRevision",
			want:    false,
		},
		{
			name:     "not_revision_bad_name",
			Message:  "Book",
			Resource: `option (google.api.resource) = {type: "library.googleapis.com/Book"};`,
			want:     false,
		},
	} {
		f := testutils.Compile(t, `
			import "google/api/resource.proto";
			message {{.Message}} {
				{{.Resource}}
				string name = 1;
			}
		`, test)
		m := f.Messages().Get(0)
		if got := IsResourceRevision(m); got != test.want {
			t.Errorf("IsResourceRevision(%+v): got %v, want %v", m, got, test.want)
		}
	}
}

func TestIsRevisionRelationship(t *testing.T) {
	for _, test := range []struct {
		name         string
		typeA, typeB string
		want         bool
	}{
		{
			name:  "revision_relationship",
			typeA: "library.googleapis.com/Book",
			typeB: "library.googleapis.com/BookRevision",
			want:  true,
		},
		{
			name:  "non_revision_relationship",
			typeA: "library.googleapis.com/Book",
			typeB: "library.googleapis.com/Library",
			want:  false,
		},
		{
			name:  "invalid_type_a",
			typeA: "library.googleapis.com",
			typeB: "library.googleapis.com/Library",
			want:  false,
		},
		{
			name:  "invalid_type_b",
			typeA: "library.googleapis.com/Book",
			typeB: "library.googleapis.com",
			want:  false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			a := &annotations.ResourceDescriptor{Type: test.typeA}
			b := &annotations.ResourceDescriptor{Type: test.typeB}
			if got := IsRevisionRelationship(a, b); got != test.want {
				t.Errorf("IsRevisionRelationship(%s, %s): got %v, want %v", test.typeA, test.typeB, got, test.want)
			}
		})
	}
}

func TestHasParent(t *testing.T) {
	for _, test := range []struct {
		name    string
		pattern string
		want    bool
	}{
		{
			name:    "child resource",
			pattern: "foos/{foo}/bars/{bar}",
			want:    true,
		},
		{
			name:    "top level resource",
			pattern: "foos/{foo}",
			want:    false,
		},
		{
			name:    "top level singleton",
			pattern: "foo",
			want:    false,
		},
		{
			name:    "empty",
			pattern: "",
			want:    false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var in *annotations.ResourceDescriptor
			if test.pattern != "" {
				in = &annotations.ResourceDescriptor{Pattern: []string{test.pattern}}
			}
			if got := HasParent(in); got != test.want {
				t.Errorf("HasParent(%s): got %v, want %v", test.pattern, got, test.want)
			}
		})
	}
}

