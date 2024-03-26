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

package lint

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v2"
)

// Configs determine if a rule is enabled or not on a file path.
type Configs []Config

// Config stores rule configurations for certain files
// that the file path must match any of the included paths
// but none of the excluded ones.
type Config struct {
	IncludedPaths []string `json:"included_paths" yaml:"included_paths"`
	ExcludedPaths []string `json:"excluded_paths" yaml:"excluded_paths"`
	EnabledRules  []string `json:"enabled_rules" yaml:"enabled_rules"`
	DisabledRules []string `json:"disabled_rules" yaml:"disabled_rules"`
	ImportPaths   []string `json:"import_paths" yaml:"import_paths"`
}

// ReadConfigsFromFile reads Configs from a file.
// It supports JSON(.json) and YAML(.yaml or .yml) files.
func ReadConfigsFromFile(path string) (Configs, error) {
	var parse func(io.Reader) (Configs, error)
	switch filepath.Ext(path) {
	case ".json":
		parse = ReadConfigsJSON
	case ".yaml", ".yml":
		parse = ReadConfigsYAML
	}
	if parse == nil {
		return nil, fmt.Errorf("reading Configs: unsupported format `%q` with file path `%q`", filepath.Ext(path), path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readConfig: %s", err.Error())
	}
	defer f.Close()

	return parse(f)
}

// ReadConfigsJSON reads Configs from a JSON file.
func ReadConfigsJSON(f io.Reader) (Configs, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Configs
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// ReadConfigsYAML reads Configs from a YAML(.yml or .yaml) file.
func ReadConfigsYAML(f io.Reader) (Configs, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Configs
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// IsRuleEnabled returns true if a rule is enabled by the configs.
func (configs Configs) IsRuleEnabled(rule string, path string) bool {
	// Enabled by default if the rule does not belong to one of the default
	// disabled groups. Otherwise, needs to be explicitly enabled.
	enabled := !matchRule(rule, defaultDisabledRules...)
	for _, c := range configs {
		if c.matchPath(path) {
			if matchRule(rule, c.DisabledRules...) {
				enabled = false
			}
			if matchRule(rule, c.EnabledRules...) {
				enabled = true
			}
		}
	}

	return enabled
}

func (c Config) matchPath(path string) bool {
	if matchPath(path, c.ExcludedPaths...) {
		return false
	}
	return len(c.IncludedPaths) == 0 || matchPath(path, c.IncludedPaths...)
}

func matchPath(path string, pathPatterns ...string) bool {
	for _, pattern := range pathPatterns {
		if matched, _ := doublestar.Match(pattern, path); matched {
			return true
		}
	}
	return false
}

func matchRule(rule string, rulePrefixes ...string) bool {
	rule = strings.ToLower(rule)
	for _, prefix := range rulePrefixes {
		prefix = strings.ToLower(prefix)
		prefix = strings.TrimSuffix(prefix, nameSeparator) // "core::" -> "core"
		prefix = strings.TrimPrefix(prefix, nameSeparator) // "::http-body" -> "http-body"
		if prefix == "all" ||
			prefix == rule ||
			strings.HasPrefix(rule, prefix+nameSeparator) || // e.g., "core" matches "core::http-body", but not "core-rules::http-body"
			strings.HasSuffix(rule, nameSeparator+prefix) || // e.g., "http-body" matches "core::http-body", but not "core::google-http-body"
			strings.Contains(rule, nameSeparator+prefix+nameSeparator) { // e.g., "http-body" matches "core::http-body::post", but not "core::google-http-body::post"
			return true
		}
	}
	return false
}
