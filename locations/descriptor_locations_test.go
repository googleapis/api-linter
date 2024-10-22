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

func TestDescriptorName(t *testing.T) {
	f := parse(t, `
		message Foo {
		  string bar = 1;
		  map<string, string> baz = 2;
		}
	`)

	tests := []struct {
		testName string
		d        desc.Descriptor
		wantSpan []int32
	}{
		{"Message", f.GetMessageTypes()[0], []int32{2, 8, 11}},
		{"Field", f.GetMessageTypes()[0].GetFields()[0], []int32{3, 9, 12}},
		{"MapField", f.GetMessageTypes()[0].GetFields()[1], []int32{4, 22, 25}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if diff := cmp.Diff(DescriptorName(test.d).Span, test.wantSpan); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestFileDescriptorName(t *testing.T) {
	f := parse(t, `
		message Foo {}
	`)
	if got := DescriptorName(f); got != nil {
		t.Errorf("%v", got)
	}
}
