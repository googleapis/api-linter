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

import "google.golang.org/protobuf/reflect/protoreflect"

// Response describes the result returned by a rule.
type Response struct {
	FilePath string    `json:"file_path" yaml:"file_path"`
	Problems []Problem `json:"problems" yaml:"problems"`
}

// Problem contains information about a result produced by an API Linter.
type Problem struct {
	// Message provides a short description of the problem.
	Message string `json:"message" yaml:"message"`
	// Suggestion provides a suggested fix, if applicable.
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
	// Location provides the location of the problem. If both
	// `Location` and `Descriptor` are specified, the location
	// is then used from `Location` instead of `Descriptor`.
	Location Location `json:"location" yaml:"location"`
	// Descriptor provides the descriptor related
	// to the problem. If present and `Location` is not
	// specified, then the starting location of the descriptor
	// is used as the location of the problem.
	Descriptor protoreflect.Descriptor `json:"-" yaml:"-"`

	// RuleID provides the ID of the rule that this problem belongs to.
	// DO NOT SET: this field will be set by the linter based on rule info
	// and user configs.
	RuleID RuleName `json:"rule_id" yaml:"rule_id"`

	// DO NOT SET:  this field will be set by the linter based on user configs.
	Category string `json:"category,omitempty" yaml:"category,omitempty"`
}
