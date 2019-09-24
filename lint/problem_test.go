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

package lint

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestEqualTrue(t *testing.T) {
	// Build a message for the descriptor test.
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Declare a series of tests that should all be equal.
	tests := []struct {
		name string
		x    Problems
		y    Problems
	}{
		{"NilNil", nil, nil},
		{"ProblemNil", Problems{}, nil},
		{"Descriptor", Problems{{Descriptor: m}}, Problems{{Descriptor: m}}},
		{"Suggestion", Problems{{Suggestion: "foo"}}, Problems{{Suggestion: "foo"}}},
		{"MessageExact", Problems{{Message: "foo"}}, Problems{{Message: "foo"}}},
		{"MessageSubstr", Problems{{Message: "foo"}}, Problems{{Message: "foo bar"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.x, test.y); diff != "" {
				t.Errorf("Problems were unequal (x, y):\n%v", diff)
			}
			if diff := cmp.Diff(test.y, test.x); diff != "" {
				t.Errorf("Problems were unequal (y, x):\n%v", diff)
			}
		})
	}
}

func TestEqualFalse(t *testing.T) {
	// Build a message for the descriptor test.
	m1, err1 := builder.NewMessage("Foo").Build()
	m2, err2 := builder.NewMessage("Bar").Build()
	if err1 != nil || err2 != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Declare a series of tests that should all be equal.
	tests := []struct {
		name string
		x    Problems
		y    Problems
	}{
		{"ProblemNil", Problems{{}}, nil},
		{"Descriptor", Problems{{Descriptor: m1}}, Problems{{Descriptor: m2}}},
		{"Suggestion", Problems{{Suggestion: "foo"}}, Problems{{Suggestion: "bar"}}},
		{"Message", Problems{{Message: "foo"}}, Problems{{Message: "bar"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.x, test.y); diff == "" {
				t.Errorf("Got no diff (x, y); expected one.")
			}
			if diff := cmp.Diff(test.y, test.x); diff == "" {
				t.Errorf("Got no diff (y, x); expected one.")
			}
		})
	}
}

func TestSetDescriptor(t *testing.T) {
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}
	problems := Problems{{}, {}, {}}.SetDescriptor(m)
	for _, p := range problems {
		if p.Descriptor != m {
			t.Errorf("Got %v, expected %v", p.Descriptor, m)
		}
	}
}
