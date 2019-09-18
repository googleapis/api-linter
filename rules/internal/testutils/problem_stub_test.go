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

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestVerifyDescriptor(t *testing.T) {
	// Build a descriptor to verify with.
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Unable to build descriptor.")
	}

	// Establish that the descriptor matches.
	stub := ProblemStub{Descriptor: m}
	problem := lint.Problem{Descriptor: m}
	stub.VerifyDescriptor(problem, t)
}

func TestVerifyDescriptorError(t *testing.T) {
	// Build a descriptor to verify with.
	mx, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Unable to build descriptor.")
	}
	my, err := builder.NewMessage("Bar").Build()
	if err != nil {
		t.Fatalf("Unable to build descriptor.")
	}

	// Establish that the descriptor matches.
	stub := ProblemStub{Descriptor: mx}
	problem := lint.Problem{Descriptor: my}
	canary := &testing.T{}
	stub.VerifyDescriptor(problem, canary)
	if !canary.Failed() {
		t.Errorf("ProblemStub.VerifyDescriptor succeeded; expected error.")
	}
}

func TestVerifyMessage(t *testing.T) {
	stub := ProblemStub{Message: "foo"}
	problem := lint.Problem{Message: "foo bar"}
	stub.VerifyMessage(problem, t)
}

func TestVerifyMessageError(t *testing.T) {
	canary := &testing.T{}
	stub := ProblemStub{Message: "foo"}
	problem := lint.Problem{Message: "bar baz"}
	stub.VerifyMessage(problem, canary)
	if !canary.Failed() {
		t.Errorf("ProblemStub.VerifyMessage succeeded; expected error.")
	}
}

func TestSetDescriptor(t *testing.T) {
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Unable to build descriptor.")
	}
	stubs := ProblemStubs{ProblemStub{}, ProblemStub{}}
	resp := stubs.SetDescriptor(m)
	for _, stub := range resp {
		if stub.Descriptor != m {
			t.Errorf("ProblemStubs.SetDescriptor did not set the descriptor properly.")
		}
	}
}

func TestVerifyProblemStubs(t *testing.T) {
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("Unable to build descriptor.")
	}

	// Build a ProblemStubs object.
	stubs := ProblemStubs{
		ProblemStub{Descriptor: m, Message: "foo"},
		ProblemStub{Descriptor: m},
	}

	// Build a slice of problem objects that should match.
	problems := []lint.Problem{
		lint.Problem{Descriptor: m, Message: "foo bar"},
		lint.Problem{Descriptor: m, Message: "foo bar"},
	}

	// Verify the problems against the stubs.
	stubs.Verify(problems, t)
}

func TestVerifyProblemStubsMismatchedLength(t *testing.T) {
	canary := &testing.T{}
	stubs := ProblemStubs{ProblemStub{}}
	problems := []lint.Problem{}
	stubs.Verify(problems, canary)
	if !canary.Failed() {
		t.Errorf("stubs.Verify succeeded; expected error.")
	}
}
