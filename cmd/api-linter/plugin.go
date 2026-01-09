// Copyright 2025 Google LLC
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
	"plugin"

	"github.com/googleapis/api-linter/lint"
)

// addRulesFuncName is the expected name of the function exported by plugins.
const addRulesFuncName = "AddCustomRules"

// loadCustomRulePlugin loads a plugin from the given path and registers
// its rules with the provided registry.
func loadCustomRulePlugin(pluginPath string, registry lint.RuleRegistry) error {
	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %q: %v", pluginPath, err)
	}

	// Look up the AddCustomRules function
	sym, err := p.Lookup(addRulesFuncName)
	if err != nil {
		return fmt.Errorf("plugin %q does not export %q function: %v",
			pluginPath, addRulesFuncName, err)
	}

	// Cast to the expected function type
	addRulesFunc, ok := sym.(func(lint.RuleRegistry) error)
	if !ok {
		return fmt.Errorf("plugin %q exports %q with wrong signature",
			pluginPath, addRulesFuncName)
	}

	// Call the function to add custom rules
	if err := addRulesFunc(registry); err != nil {
		return fmt.Errorf("error registering rules from plugin %q: %v",
			pluginPath, err)
	}

	return nil
}

// loadCustomRulePlugins loads all plugins from the given paths and registers
// their rules with the provided registry.
func loadCustomRulePlugins(pluginPaths []string, registry lint.RuleRegistry) error {
	for _, path := range pluginPaths {
		if err := loadCustomRulePlugin(path, registry); err != nil {
			return err
		}
	}
	return nil
}
