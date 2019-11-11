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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func checkRulesRegistered() (errors []error) {
	// Iterate over each file in the rules/aip* directories and do two things:
	//   (1) Make a list of expected AIPs.
	//   (2) Make sure each rule is registered in that AIP's `AddRules` function.
	allRules := map[int][]string{}
	err := filepath.Walk("./rules", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err)
			return nil
		}

		// If this is a directory or a test file, skip it.
		if info.IsDir() || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Read in the file.
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			errors = append(errors, err)
			return nil
		}

		// Look for a rule in the file, as well as the AIP the rule is from.
		if aipMatch, ruleMatch := aipDirRegexp.FindStringSubmatch(path), ruleRegexp.FindSubmatch(contents); len(aipMatch) > 0 && len(ruleMatch) > 0 {
			rule := string(ruleMatch[1])
			aip, err := strconv.Atoi(aipMatch[1])
			if err != nil {
				errors = append(errors, err)
				return nil
			}

			// Add the rule to our list of rules.
			if _, ok := allRules[aip]; !ok {
				allRules[aip] = []string{}
			}
			allRules[aip] = append(allRules[aip], rule)
		}
		return nil
	})
	if err != nil {
		errors = append(errors, err)
	}

	// Iterate over all of the rules we found and ensure that they are registered,
	// and that the AIP itself is registered.
	contentsBytes, err := ioutil.ReadFile("rules/rules.go")
	if err != nil {
		errors = append(errors, err)
	}
	contents := string(contentsBytes)
	for aip, rules := range allRules {
		// Make sure the AIP's package's rules are registered in coreRules in rules.go.
		if !strings.Contains(contents, fmt.Sprintf("aip%04d.AddRules(coreRules)", aip)) {
			errors = append(errors, fmt.Errorf("rules.go does not call AllRules for for AIP-%d", aip))
		}

		// Make sure each individual rule is included in AllRules.
		aipBaseFileContentsBytes, err := ioutil.ReadFile(fmt.Sprintf("rules/aip%04d/aip%04d.go", aip, aip))
		if err != nil {
			errors = append(errors, err)
		}
		aipBaseFileContents := string(aipBaseFileContentsBytes)
		for _, rule := range rules {
			if !strings.Contains(aipBaseFileContents, rule) {
				errors = append(errors, fmt.Errorf("rule %q for AIP-%d not registered in the AIP's AllRules", rule, aip))
			}
		}
	}

	return
}

var aipDirRegexp = regexp.MustCompile(`/aip([\d]{4})/`)
var ruleRegexp = regexp.MustCompile(`var ([\w\d]+) = &lint\.[\w]+Rule\s*{`)
