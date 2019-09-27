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
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Problems is a slice of individual Problem objects.
type Problems []lint.Problem

// Diff determines whether a Problem is sufficiently similar to another
// to be considered equivalent, and returns a diff otherwise.
//
// This is intended for unit tests and is intentially generous on what
// constitutes equality.
func (problems Problems) Diff(other []lint.Problem) string {
	// `other` may be nil.
	// Consider an length-0 slice to be equal to nil.
	if other == nil {
		if len(problems) == 0 {
			return ""
		}
		return cmp.Diff(problems, other)
	}

	// If the problems differ in length, they are by definition unequal.
	if len(problems) != len(other) {
		return cmp.Diff(problems, other)
	}

	// Iterate over the individual problems and determine whether they are
	// sufficiently equivalent.
	for i := range problems {
		x, y := problems[i], other[i]

		// The descriptors must exactly match, otherwise the problems are unequal.
		if x.Descriptor != y.Descriptor {
			return cmp.Diff(problems, other)
		}

		// The suggestions, if present, must exactly match.
		if x.Suggestion != y.Suggestion {
			return cmp.Diff(problems, other)
		}

		// When comparing messages, we want to know if the test string is a
		// substring of the actual one.
		if !strings.Contains(y.Message, x.Message) {
			return cmp.Diff(problems, other)
		}
	}

	// These sets of problems are sufficiently equal.
	return ""
}

// SetDescriptor sets the given descriptor to every Problem in the slice and
// returns the slice back.
//
// This is intended primarily for use in unit tests.
func (problems Problems) SetDescriptor(d desc.Descriptor) Problems {
	for i := range problems {
		problems[i].Descriptor = d
	}
	return problems
}
