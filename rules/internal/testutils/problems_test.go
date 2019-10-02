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

package testutils

import (
	"testing"

	. "github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestDiffEquivalent(t *testing.T) {
	// Build a message for the descriptor test.
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Declare a series of tests that should all be equal.
	tests := []struct {
		name string
		x    Problems
		y    []Problem
	}{
		{"NilNil", nil, nil},
		{"ProblemNil", Problems{}, nil},
		{"Descriptor", Problems{{Descriptor: m}}, []Problem{{Descriptor: m}}},
		{"Suggestion", Problems{{Suggestion: "foo"}}, []Problem{{Suggestion: "foo"}}},
		{"MessageExact", Problems{{Message: "foo"}}, []Problem{{Message: "foo"}}},
		{"MessageSubstr", Problems{{Message: "foo"}}, []Problem{{Message: "foo bar"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if diff := test.x.Diff(test.y); diff != "" {
				t.Errorf("Problems were unequal (x, y):\n%v", diff)
			}
		})
	}
}

func TestDiffNotEquivalent(t *testing.T) {
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
		y    []Problem
	}{
		{"ProblemNil", Problems{{Descriptor: m1}}, nil},
		{"EmptyProblemNil", Problems{{}}, nil},
		{"LengthMismatch", Problems{{}}, []Problem{{}, {}}},
		{"Descriptor", Problems{{Descriptor: m1}}, []Problem{{Descriptor: m2}}},
		{"Suggestion", Problems{{Suggestion: "foo"}}, []Problem{{Suggestion: "bar"}}},
		{"Message", Problems{{Message: "foo"}}, []Problem{{Message: "bar"}}},
		{"MessageSuperstr", Problems{{Message: "foo bar"}}, []Problem{{Message: "foo"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if diff := test.x.Diff(test.y); diff == "" {
				t.Errorf("Got no diff (x, y); expected one.")
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
