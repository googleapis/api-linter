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

package aip0140

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestAddRules(t *testing.T) {
	rules := make(lint.RuleRegistry)
	AddRules(rules)
	for ruleName := range rules {
		if !strings.HasPrefix(string(ruleName), "core::0140") {
			t.Errorf("Rule %s is not namespaced to core::0140.", ruleName)
		}
	}
}
