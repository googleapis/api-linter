// Copyright 2019 Google LLC
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

package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/stoewer/go-strcase"
)

func checkRuleRegistered(aip int, name string) []error {
	path := fmt.Sprintf("rules/aip%04d/%s.go", aip, strcase.SnakeCase(name))

	// Read in the file.
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return []error{err}
	}

	// Look for a rule in the file. Complain if we can not find one.
	ruleMatch := ruleRegexp.FindStringSubmatch(string(contents))
	if len(ruleMatch) == 0 {
		return []error{fmt.Errorf("no rule found: %s", path)}
	}

	// Some errors are now non-fatal; start a running tab.
	errata := []error{}

	// Ensure this rule is registered within its AIP module file.
	ruleVar := ruleMatch[1]
	contents, err = ioutil.ReadFile(fmt.Sprintf("rules/aip%04d/aip%04d.go", aip, aip))
	if err != nil {
		return []error{err}
	}
	if !strings.Contains(string(contents), ruleVar) {
		errata = append(errata, fmt.Errorf("rule %q for AIP-%d not registered in the AIP's AllRules", name, aip))
	}

	// Ensure that the AIP itself is registered in `rules/rules.go`.
	contents, err = ioutil.ReadFile("rules/rules.go")
	if err != nil {
		errata = append(errata, err)
		return errata
	}
	if !strings.Contains(string(contents), fmt.Sprintf("aip%04d.AddRules", aip)) {
		errata = append(errata, fmt.Errorf("rules.go does not call AllRules for for AIP-%d", aip))
	}

	// Done; return any errata we found.
	return errata
}

var ruleRegexp = regexp.MustCompile(`var ([\w\d]+) = &lint\.[\w]+Rule\s*{`)
