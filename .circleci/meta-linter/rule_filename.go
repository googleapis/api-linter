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
	"strings"
)

func ruleFilename() (errors []error) {
	// Iterate over each file in the rules directory.
	err := filepath.Walk("./rules", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err)
			return nil
		}

		// If this is a directory, skip it.
		if info.IsDir() {
			return nil
		}

		// Special case: The comment in rules/rules.go trips this check.
		if path == "rules/rules.go" {
			return nil
		}

		// Read in the file.
		contentsBytes, err := ioutil.ReadFile(path)
		if err != nil {
			errors = append(errors, err)
		}
		contents := string(contentsBytes)

		// If the file creates a rule that does not match the filename, complain.
		for _, match := range ruleNameRegexp.FindAllStringSubmatch(contents, -1) {
			aip, ruleName := match[1], match[2]
			if !strings.Contains(path, aip) {
				errors = append(errors, fmt.Errorf("mismatch between path and rule AIP: %s", path))
			}
			if !strings.HasSuffix(path, strings.ReplaceAll(ruleName, "-", "_")+".go") {
				errors = append(errors, fmt.Errorf("mismatch between rule name and filename: %s", path))
			}
		}

		return nil
	})
	if err != nil {
		errors = append(errors, err)
	}
	return
}

var ruleNameRegexp = regexp.MustCompile(`NewRuleName\(([\d]+), "([a-z-]+)"\)`)
