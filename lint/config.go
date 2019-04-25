package lint

import (
	"encoding/json"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"io"
	"io/ioutil"
)

// RuntimeConfigs stores a list of RuntimeConfig and supports config lookup for a given path.
type RuntimeConfigs []RuntimeConfig

// RuntimeConfig stores rule runtime configurations and file spec that  a path must match any of the
// included paths but none of the excluded paths.
type RuntimeConfig struct {
	IncludedPaths []string              `json:"included_paths"`
	ExcludedPaths []string              `json:"excluded_paths"`
	RuleConfigs   map[string]RuleConfig `json:"rule_configs"`
}

// RuleConfig stores the status and category of a rule, which can be applied to a rule during runtime.
type RuleConfig struct {
	Status   Status   `json:"status"`
	Category Category `json:"category"`
}

var defaultRuleConfig = RuleConfig{Status: Disabled, Category: Warning}

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

// getRuleConfig returns a RuleConfig matching path and rule
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
		if matched, err := doublestar.PathMatch(pattern, path); matched || err != nil {
			return false
		}
	}
	for _, pattern := range c.IncludedPaths {
		if matched, err := doublestar.PathMatch(pattern, path); matched && err == nil {
			return true
		}
	}
	return false
}
