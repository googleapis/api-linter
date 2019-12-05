// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package locations

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc"
)

func TestFieldLocations(t *testing.T) {
	f := parse(t, `
		message Book {
		  string name = 1;
		  Author author = 2;
		}
		message Author {}
	`)
	tests := []struct {
		name  string
		field *desc.FieldDescriptor
		span  []int32
	}{
		{"Primitive", f.GetMessageTypes()[0].GetFields()[0], []int32{3, 2, 8}},
		{"Composite", f.GetMessageTypes()[0].GetFields()[1], []int32{4, 2, 8}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := FieldType(test.field)
			if diff := cmp.Diff(l.GetSpan(), test.span); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFieldResourceReference(t *testing.T) {
	f := parse(t, `
		import "google/api/resource.proto";
		message GetBookRequest {
		  string name = 1 [(google.api.resource_reference) = {
		    type: "library.googleapis.com/Book"
		  }];
		}
	`)
	loc := FieldResourceReference(f.GetMessageTypes()[0].GetFields()[0])
	if diff := cmp.Diff(loc.GetSpan(), []int32{4, 19, 6, 3}); diff != "" {
		t.Errorf(diff)
	}
}
