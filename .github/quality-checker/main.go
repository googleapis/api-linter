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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"bitbucket.org/creachadair/stringset"
)

// The set of checkers that are run on every discovered rule.
var checkers = []func(aip int, name string) []error{
	checkRuleDocumented,
	checkRuleName,
	checkRuleRegistered,
}

func main() {
	errors := []error{}

	// Keep track of rules processed, and which ones have one or more failures.
	passedRules := stringset.New()
	failedRules := stringset.New()

	// Begin walking the rules directory looking for rules.
	err := filepath.Walk("./rules", func(path string, info os.FileInfo, err error) error {
		// Sanity check: Bubble errors.
		if err != nil {
			errors = append(errors, err)
			return nil
		}

		// Weed out files that are not rules (tests, etc.).
		for _, exempt := range []func(path string, info os.FileInfo) bool{
			func(_ string, i os.FileInfo) bool { return i.IsDir() },
			func(p string, _ os.FileInfo) bool { return strings.Contains(p, "/internal/") },
			func(p string, _ os.FileInfo) bool { return p == "rules/rules.go" },
			func(p string, _ os.FileInfo) bool { return strings.HasSuffix(p, "_test.go") },
			func(p string, _ os.FileInfo) bool { return aipIndex.MatchString(p) },
		} {
			if exempt(path, info) {
				return nil
			}
		}

		// This represents a rule. Get the AIP and rule name from the path.
		match := ruleFile.FindStringSubmatch(path)
		if match == nil {
			errors = append(errors, fmt.Errorf("unexpected path: %s", path))
			return nil
		}

		// Get the AIP number and final rule segment.
		aip, err := strconv.Atoi(match[1])
		if err != nil {
			errors = append(errors, err)
			return nil
		}
		name := strings.ReplaceAll(match[2], "_", "-")
		token := fmt.Sprintf("%04d-%s", aip, name)

		// Run each checker and run up the list of errors.
		for _, checker := range checkers {
			if errs := checker(aip, name); len(errs) > 0 {
				errors = append(errors, errs...)
				failedRules.Add(token)
			}
		}

		// All checkers are done; add this to the success list if nothing failed.
		if !failedRules.Contains(token) {
			passedRules.Add(token)
		}

		return nil
	})
	// Ensure the rollup error is nil. (It should be, since our walk function
	// never returns an error but always appends instead.)
	if err != nil {
		errors = append(errors, err)
	}

	// If we got complaints, complain about them.
	if len(errors) > 0 {
		for _, e := range errors {
			log.Println(fmt.Sprintf("ERROR: %s", e.Error()))
		}
	}

	// Provide a summary.
	fmt.Printf(
		"%d rules scanned: %d passed, %d failed.\n",
		len(passedRules)+len(failedRules),
		len(passedRules),
		len(failedRules),
	)

	// Exit.
	if len(errors) > 0 {
		os.Exit(1)
	}
}

var (
	ruleFile = regexp.MustCompile(`rules/aip([\d]{4})/([a-z0-9_]+).go`)
	aipIndex = regexp.MustCompile(`rules/aip[\d]{4}/aip[\d]{4}\.go`)
)
