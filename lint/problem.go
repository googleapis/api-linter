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
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// Problem contains information about a result produced by an API Linter.
type Problem struct {
	// Message provides a short description of the problem.
	Message string `json:"message" yaml:"message"`

	// Suggestion provides a suggested fix, if applicable.
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`

	// Descriptor provides the descriptor related
	// to the problem. If present and `Location` is not
	// specified, then the starting location of the descriptor
	// is used as the location of the problem.
	Descriptor desc.Descriptor `json:"-" yaml:"-"`

	// Location provides the location of the problem.
	// DO NOT SET: Set the descriptor instead.
	Location Location `json:"location" yaml:"location"`

	// RuleID provides the ID of the rule that this problem belongs to.
	// DO NOT SET: this field will be set by the linter based on rule info
	// and user configs.
	RuleID RuleName `json:"rule_id" yaml:"rule_id"`

	// DO NOT SET:  this field will be set by the linter based on user configs.
	Category string `json:"category,omitempty" yaml:"category,omitempty"`

	noPositional struct{}
}

// Problems is a slice of individual Problem objects.
type Problems []Problem

// Equal determines whether a Problem is sufficiently similar to another
// to be considered equal.
//
// This is intended for unit tests and is intentially generous on what
// constitutes equality.
func (problems Problems) Equal(other Problems) bool {
	// `other` may be nil.
	// Consider an length-0 slice to be equal to nil.
	if other == nil {
		return len(problems) == 0
	}

	// If the problems differ in length, they are by definition unequal.
	if len(problems) != len(other) {
		return false
	}

	// Iterate over the individual problems and determine whether they are
	// sufficiently equivalent.
	for i := range problems {
		x, y := problems[i], other[i]

		// The descriptors must exactly match, otherwise the problems are unequal.
		if x.Descriptor != y.Descriptor {
			return false
		}

		// The suggestions, if present, must exactly match.
		if x.Suggestion != y.Suggestion {
			return false
		}

		// When comparing messages, we want to know if the shorter string is a
		// substring of the longer one, but to preserve the communitive property
		// of equality, we do not care which is which.
		short, long := x.Message, y.Message
		if len(long) < len(short) {
			short, long = long, short
		}
		if !strings.Contains(long, short) {
			return false
		}
	}

	// These sets of problems are sufficiently equal.
	return true
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
