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

// Package rules contains implementations of rules that apply to all Google APIs.
package rules

import (
	"log"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/aip0131"
	"github.com/googleapis/api-linter/rules/aip0132"
	"github.com/googleapis/api-linter/rules/aip0140"
	"github.com/googleapis/api-linter/rules/aip0191"
)

func init() {
	aip0131.AddRules(coreRules)
	aip0132.AddRules(coreRules)
	aip0140.AddRules(coreRules)
	aip0191.AddRules(coreRules)
}

var coreRules, _ = lint.NewRules()

// Rules returns all rules registered in this package.
func Rules() lint.Rules {
	return coreRules.Copy()
}

// registerRules registers the given rule into "core rules".
func registerRules(r ...lint.Rule) {
	for _, rl := range r {
		if err := coreRules.Register(rl); err != nil {
			log.Fatalf("Error when registering rule '%s': %v", rl.GetName(), err)
		}
	}
}
