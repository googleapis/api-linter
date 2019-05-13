package lint

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	"gopkg.in/yaml.v2"
)

// Configs stores a list of Config.
//
// Note: for a path, if multiple configs match it, the rule configs of the later one
// will always override those in the former one. For example, given a list of configs
//
// [
//		{
//			included_paths: ["/a/**/*.proto"],
//			rule_configs: {"my::rule": {disabled, error}},
//		},
//		{
//			included_paths: ["/a/b.proto"],
//			rule_configs: {"my::rule": {enabled}},
//		}
// ]
//
// that match path "/a/b.proto", the resulted config for rule "my::rule" will be {enabled, error}.
type Configs []Config

// Config stores rule configurations for certain files
// that the file path must match any of the included paths
// but none of the excluded ones.
type Config struct {
	IncludedPaths []string              `json:"included_paths" yaml:"included_paths"`
	ExcludedPaths []string              `json:"excluded_paths" yaml:"excluded_paths"`
	RuleConfigs   map[string]RuleConfig `json:"rule_configs" yaml:"rule_configs"`
}

// RuleConfig stores configurable status and category of a rule.
type RuleConfig struct {
	Disabled bool   `json:"disabled" yaml:"disabled"`
	Category string `json:"category" yaml:"category"`
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
		return nil, fmt.Errorf("Reading Configs: unsupported format `%q` with file path `%q`", filepath.Ext(path), path)
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
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Configs
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// ReadConfigsYAML reads Configs from a JSON file.
func ReadConfigsYAML(f io.Reader) (Configs, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Configs
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// getRuleConfig returns a RuleConfig that matches the given path and rule.
// Returns an error if a config is not found for the path.
func (c Configs) getRuleConfig(path string, rule RuleName) (result RuleConfig, err error) {
	err = fmt.Errorf("failed to find a config for path %q", path)
	for _, cfg := range c {
		if cfg.match(path) {
			err = nil
			if r, ok := cfg.getRuleConfig(rule); ok {
				result = result.withOverride(r)
			}
		}
	}
	return
}

func (c Config) getRuleConfig(rule RuleName) (RuleConfig, bool) {
	for r := rule; ; r = r.parent() {
		if ruleConfig, ok := c.RuleConfigs[string(r)]; ok {
			return ruleConfig, true
		}

		if r == "" {
			break
		}
	}

	return RuleConfig{}, false
}

// withOverride returns a copy of r, overridden with non-zero values in r2
func (r RuleConfig) withOverride(r2 RuleConfig) RuleConfig {
	r.Disabled = r2.Disabled

	if r2.Category != "" {
		r.Category = r2.Category
	}

	return r
}

// match returns if a Config matches path based on its included and excluded paths
func (c Config) match(path string) bool {
	for _, pattern := range c.ExcludedPaths {
		if matched, err := doublestar.Match(pattern, path); matched || err != nil {
			return false
		}
	}
	for _, pattern := range c.IncludedPaths {
		if matched, err := doublestar.Match(pattern, path); matched && err == nil {
			return true
		}
	}
	return false
}
