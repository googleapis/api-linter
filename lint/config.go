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

// RuntimeConfigs stores a list of RuntimeConfig.
//
// Note: for a path, if multiple runtime configs match it, the rule configs of the later one
// will always override those in the former one. For example, given a list of runtime configs
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
type RuntimeConfigs []RuntimeConfig

// RuntimeConfig stores rule configurations for certain files
// that the file path must match any of the included paths
// but none of the excluded ones.
type RuntimeConfig struct {
	IncludedPaths []string              `json:"included_paths" yaml:"included_paths"`
	ExcludedPaths []string              `json:"excluded_paths" yaml:"excluded_paths"`
	RuleConfigs   map[string]RuleConfig `json:"rule_configs" yaml:"rule_configs"`
}

// RuleConfig stores runtime-configurable status and category of a rule.
type RuleConfig struct {
	Status   Status   `json:"status" yaml:"status"`
	Category Category `json:"category" yaml:"category"`
}

func getDefaultRuleConfig() RuleConfig {
	return RuleConfig{Status: Disabled, Category: Warning}
}

// ReadConfigsFromFile reads RuntimeConfigs from a file.
// It supports JSON(.json) and YAML(.yaml or .yml) files.
func ReadConfigsFromFile(path string) (RuntimeConfigs, error) {
	var parse func(io.Reader) (RuntimeConfigs, error)
	switch filepath.Ext(path) {
	case ".json":
		parse = ReadConfigsJSON
	case ".yaml", ".yml":
		parse = ReadConfigsYAML
	}
	if parse == nil {
		return nil, fmt.Errorf("Reading RuntimeConfigs: unsupported format `%q` with file path `%q`", filepath.Ext(path), path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readConfig: %s", err.Error())
	}
	defer f.Close()

	return parse(f)
}

// ReadConfigsJSON reads RuntimeConfigs from a JSON file.
func ReadConfigsJSON(f io.Reader) (RuntimeConfigs, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c RuntimeConfigs
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// ReadConfigsYAML reads RuntimeConfigs from a JSON file.
func ReadConfigsYAML(f io.Reader) (RuntimeConfigs, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c RuntimeConfigs
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// getRuleConfig returns a RuleConfig that matches the given path and rule.
// Returns an error if a config is not found for the path.
func (c RuntimeConfigs) getRuleConfig(path string, rule RuleName) (result RuleConfig, err error) {
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

func (c RuntimeConfig) getRuleConfig(rule RuleName) (RuleConfig, bool) {
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
	if r2.Status != "" {
		r.Status = r2.Status
	}

	if r2.Category != "" {
		r.Category = r2.Category
	}

	return r
}

// match returns if a RuntimeConfig matches path based on its included and excluded paths
func (c RuntimeConfig) match(path string) bool {
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
