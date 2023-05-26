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

	apb "google.golang.org/genproto/googleapis/api/annotations"
)

func TestMatches(t *testing.T) {
	for _, test := range []struct {
		name     string
		resource *apb.ResourceDescriptor
		want     string
	}{
		{
			name: "SingularSpecified",
			resource: &apb.ResourceDescriptor{
				Singular: "bookShelf",
			},
			want: "bookShelf",
		},
		{
			name: "SingularAndTypeSpecified",
			resource: &apb.ResourceDescriptor{
				Singular: "bookShelf",
				// NOTE: this is not a correct resource annotation.
				// it must match singular.
				Type: "library.googleapis.com/book",
			},
			want: "bookShelf",
		},
		{
			name: "TypeSpecified",
			resource: &apb.ResourceDescriptor{
				Type: "library.googleapis.com/bookShelf",
			},
			want: "bookShelf",
		},
		{
			name:     "NothingSpecified",
			resource: &apb.ResourceDescriptor{},
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
