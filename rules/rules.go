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
// A rule is technically anything with a `GetName()`, `GetURI()“, and
// `Lint(*desc.FileDescriptorProto) []lint.Problem` method, but most rule
// authors will want to use the rule structs provided in the lint package
// (`&lint.MessageRule`, `&lint.FieldRule`, and so on). These run against
// each applicable descriptor in the file (`MessageRule` against every message,
// for example). They also have an `OnlyIf` property that can be used to run
// against a subset of descriptors.
//
// A simple rule therefore looks like this:
//
//	var myRule = &lint.MessageRule{
//	  Name: lint.NewRuleName(1234, "my-rule"),
//	  LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
//	    if isBad(m) {
//	      return []lint.Problem{{
//	        Message: "This message is bad.",
//	        Descriptor: m,
//	      }}
//	    }
//	    return nil
//	  },
//	}
//
// Once a rule is written, it needs to be registered. This involves adding
// the rule to the `AddRules` method for the appropriate AIP package.
// If this is the first rule for a new AIP, then the `rules.go` init() function
// must also be updated to run the `AddRules` function for the new package.
package rules

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/aip0121"
	// "github.com/googleapis/api-linter/rules/aip0122" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0123" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0124" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0126" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0127" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0128" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0131" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0132" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0133" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0134" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0135" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0136" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0140" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0141" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0142" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0143" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0144" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0146" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0148" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0151" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0152" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0154" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0155" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0156" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0157" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0158" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0159" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0162" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0163" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0164" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0165" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0191" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0192" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0202" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0203" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0214" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0215" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0216" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0217" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0231" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0233" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0234" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip0235" //TODO: uncommnet line and migrate rule to new interface
	// "github.com/googleapis/api-linter/rules/aip4232" //TODO: uncommnet line and migrate rule to new interface
)

type addRulesFuncType func(lint.RuleRegistry) error

var aipAddRulesFuncs = []addRulesFuncType{
	aip0121.AddRules,
	// aip0122.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0123.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0124.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0126.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0127.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0128.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0131.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0132.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0133.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0134.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0135.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0136.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0140.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0141.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0142.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0143.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0144.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0146.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0148.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0151.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0152.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0154.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0155.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0156.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0157.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0158.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0159.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0162.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0163.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0164.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0165.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0191.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0192.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0202.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0203.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0214.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0215.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0216.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0217.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0231.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0233.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0234.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip0235.AddRules, //TODO: uncommnet line and migrate rule to new interface
	// aip4232.AddRules, //TODO: uncommnet line and migrate rule to new interface
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
