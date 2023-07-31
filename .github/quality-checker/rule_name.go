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
	"os"
	"regexp"

	"github.com/stoewer/go-strcase"
)

func checkRuleName(aip int, name string) []error {
	path := fmt.Sprintf("rules/aip%04d/%s.go", aip, strcase.SnakeCase(name))

	// Read in the file.
	contentsBytes, err := os.ReadFile(path)
	if err != nil {
		return []error{err}
	}
	contents := string(contentsBytes)

	// Find the rule name declaration.
	// If it can not be found, complain.
	match := ruleNameRegexp.FindStringSubmatch(contents)
	if match == nil {
		return []error{fmt.Errorf("no rule name found: AIP-%d, %s", aip, name)}
	}

	// If the rule name declaration does not match, complain.
	errs := []error{}
	if fmt.Sprintf("%d", aip) != match[1] {
		errs = append(errs, fmt.Errorf("mismatch between path and rule AIP: %s", path))
	}
	if name != match[2] {
		errs = append(errs, fmt.Errorf("mismatch between rule name and filename: %s", path))
	}
	return errs
}

var ruleNameRegexp = regexp.MustCompile(`NewRuleName\(([\d]+), "([a-z0-9-]+)"\)`)
