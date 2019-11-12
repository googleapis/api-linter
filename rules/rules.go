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

// Package rules contains implementations of rules that apply to Google APIs.
//
// Rules are sorted into subpackages by the AIP (https://aip.dev/) that
// mandates the rule. Every rule represented in code here must be represented
// in English in a corresponding AIP. Conversely, anything mandated in an AIP
// should have a rule here if it is feasible to enforce in code (sometimes it
// is infeasible, however).
//
// A rule is technically anything with a `GetName()`, `GetURI()``, and
// `Lint(*desc.FileDescriptorProto) []lint.Problem` method, but most rule
// authors will want to use the rule structs provided in the lint package
// (`&lint.MessageRule`, `&lint.FieldRule`, and so on). These run against
// each applicable descriptor in the file (`MessageRule` against every message,
// for example). They also have an `OnlyIf` property that can be used to run
// against a subset of descriptors.
//
// A simple rule therefore looks like this:
//
//   var myRule = &lint.MessageRule{
//     Name: lint.NewRuleName(1234, "my-rule"),
//     LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
//       if isBad(m) {
//         return []lint.Problem{{
//           Message: "This message is bad.",
//           Descriptor: m,
//         }}
//       }
//       return nil
//     },
//   }
//
// Once a rule is written, it needs to be registered. This involves adding
// the rule to the `AddRules` method for the appropriate AIP package.
// If this is the first rule for a new AIP, then the `rules.go` init() function
// must also be updated to run the `AddRules` function for the new package.
package rules

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/aip0122"
	"github.com/googleapis/api-linter/rules/aip0126"
	"github.com/googleapis/api-linter/rules/aip0131"
	"github.com/googleapis/api-linter/rules/aip0132"
	"github.com/googleapis/api-linter/rules/aip0133"
	"github.com/googleapis/api-linter/rules/aip0134"
	"github.com/googleapis/api-linter/rules/aip0135"
	"github.com/googleapis/api-linter/rules/aip0136"
	"github.com/googleapis/api-linter/rules/aip0140"
	"github.com/googleapis/api-linter/rules/aip0141"
	"github.com/googleapis/api-linter/rules/aip0142"
	"github.com/googleapis/api-linter/rules/aip0143"
	"github.com/googleapis/api-linter/rules/aip0151"
	"github.com/googleapis/api-linter/rules/aip0158"
	"github.com/googleapis/api-linter/rules/aip0191"
	"github.com/googleapis/api-linter/rules/aip0192"
	"github.com/googleapis/api-linter/rules/aip0203"
	"github.com/googleapis/api-linter/rules/aip0231"
	"github.com/googleapis/api-linter/rules/aip0233"
)

type addRulesFuncType func(lint.RuleRegistry) error

var aipAddRulesFuncs = []addRulesFuncType{
	aip0122.AddRules,
	aip0126.AddRules,
	aip0131.AddRules,
	aip0132.AddRules,
	aip0133.AddRules,
	aip0134.AddRules,
	aip0135.AddRules,
	aip0136.AddRules,
	aip0140.AddRules,
	aip0141.AddRules,
	aip0142.AddRules,
	aip0143.AddRules,
	aip0151.AddRules,
	aip0158.AddRules,
	aip0191.AddRules,
	aip0192.AddRules,
	aip0203.AddRules,
	aip0231.AddRules,
  	api0233.AddRules,
}

// Add all rules to the given registry.
func Add(r lint.RuleRegistry) error {
	return addAIPRules(r, aipAddRulesFuncs)
}

func addAIPRules(r lint.RuleRegistry, addRulesFuncs []addRulesFuncType) error {
	for _, addRules := range addRulesFuncs {
		if err := addRules(r); err != nil {
			return err
		}
	}
	return nil
}
