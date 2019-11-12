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

func checkRulesDocumented() (errors []error) {
	// Iterate over each file in the rules directory.
	err := filepath.Walk("./rules", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err)
			return nil
		}

		// If this is not a rule file, skip it.
		match := ruleFile.FindStringSubmatch(path)
		if match == nil || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Determine the expected path for documentation.
		wantFile := fmt.Sprintf("docs/rules/core/%s-%s.md", match[1], strings.ReplaceAll(match[2], "_", "-"))
		if _, err := ioutil.ReadFile(wantFile); err != nil {
			errors = append(errors, fmt.Errorf("missing rule documentation: %s", wantFile))
			return nil
		}

		return nil
	})
	if err != nil {
		errors = append(errors, err)
	}
	return
}

var ruleFile = regexp.MustCompile(`rules/aip([\d]{4})/([a-z_]+).go`)
