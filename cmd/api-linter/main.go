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

// The command line `api-lint` checks Google APIs defined in Protobuf files.
// It follows the API Improvement Proposals defined in https://aip.dev.
package main

import (
	"log"
	"os"

	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
)

var rules, _ = lint.NewRuleRegistry()
var configs lint.Configs

func init() {
	configs = lint.Configs{
		lint.Config{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"core": {},
			},
		},
	}

	if err := addRules(core.Rules()); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	if err := runCLI(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}

func runCLI(args []string) error {
	c := newCli(args)
	return c.lint(rules, configs)
}

func addRules(r lint.RuleRegistry) error {
	return rules.Register(r.All()...)
}
