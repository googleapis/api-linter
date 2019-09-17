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

package testutil

import (
	"reflect"
	"strings"
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// ProblemStub describes a subset of a Problem, and is used to verify
// that problems look correct in tests.
type ProblemStub struct {
	Descriptor desc.Descriptor
	Message    string
}

// VerifyDescriptor verifies that the stub's descriptor and
// problem's descriptor match.
func (ps *ProblemStub) VerifyDescriptor(p lint.Problem, t *testing.T) {
	if !reflect.DeepEqual(p.Descriptor, ps.Descriptor) {
		t.Errorf("Expected descriptor %v; got %v", ps.Descriptor, p.Descriptor)
	}
}

// VerifyMessage verifies that the stub's message is included in the
// problem's message.
func (ps *ProblemStub) VerifyMessage(p lint.Problem, t *testing.T) {
	if !strings.Contains(p.Message, ps.Message) {
		t.Errorf(
			"Expected the problem's message to contain %q.\nActual: %q",
			ps.Message,
			p.Message,
		)
	}
}

// ProblemStubs is a slice of ProblemStub objects.
type ProblemStubs []ProblemStub

// Verify establishes that the provided problems match the stubs.
func (pss ProblemStubs) Verify(problems []lint.Problem, t *testing.T) {
	// Ensure that we got the same number of problems.
	// If we did not, then it is probably difficult to compare beyond that.
	if got, want := len(problems), len(pss); got != want {
		t.Errorf("Expected %d problems; got %d.", want, got)
		return
	}

	// Compare each of the problem stubs to establish that they
	// match the actual problems.
	for i, stub := range pss {
		stub.VerifyDescriptor(problems[i], t)
		if stub.Message != "" {
			stub.VerifyMessage(problems[i], t)
		}
	}
}
