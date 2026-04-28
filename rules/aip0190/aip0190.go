// Copyright 2024 Google LLC
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

package aip0190

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		190,
		messageCase,
		methodCase,
		serviceCase,
	)
}

var termsWithSingleLetterWords = []string{
	"AStar", "BTree", "RTree", "MTree", "QTree",
	"NGram", "NAry", "NWay", "KWay", "QLearning",
	"XAxis", "YAxis", "ZAxis", "KMeans", "PValue",
	"TTest", "TValue", "ZScore", "FTest", "FScore",
	"HIndex", "INode", "VNode", "LValue", "RValue",
	"PCode", "SExpression", "ETag", "IFrame", "VCard",
	"CName", "QName", "OAuth", "ABTest", "DPad",
	"ZIndex", "ZOrder", "ZBuffer", "XRay",
}

// validNameRegex validates that a string is a valid UpperCamelCase name
// without illegal consecutive capital letters, unless they match allowlisted terms.
var validNameRegex = regexp.MustCompile(`^(?:` +
	// A capitalized word from the prose form of the name.
	`[A-Z][a-z0-9]+|` +
	// Or an allowlisted term containing single letter "words".
	`(?:` + strings.Join(termsWithSingleLetterWords, "|") + `)[a-z0-9]*` +
`)+$`)

// isValidCamelCase checks if the given name follows AIP-190 case requirements.
func isValidCamelCase(name string) bool {
	return validNameRegex.MatchString(name)
}
