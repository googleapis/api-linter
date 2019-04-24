package lint

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

// RuntimeConfigs stores a list of RuntimeConfig and supports config lookup
// for a given path.
type RuntimeConfigs []RuntimeConfig

// Search returns the first found config that matches the path
// or an error if not found.
func (c RuntimeConfigs) Search(path string) (*RuntimeConfig, error) {
	for _, cfg := range c {
		if cfg.match(path) {
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("no config matches path %q", path)
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

// RuleConfig stores the status and category of a rule,
// which can be applied to a rule during runtime.
type RuleConfig struct {
	Status   Status   `json:"status"`
	Category Category `json:"category"`
}

// WithOverride returns a copy of r, overridden with non-zero values in r2
func (r RuleConfig) WithOverride(r2 RuleConfig) RuleConfig {
	if r2.Status != "" {
		r.Status = r2.Status
	}

	if r.Category != "" {
		r.Category = r2.Category
	}

	return r
}

// RuntimeConfig stores rule runtime configurations and file spec that
// a path must match any of the included paths but none of
// the excluded paths.
type RuntimeConfig struct {
	IncludedPaths []string              `json:"included_paths"`
	ExcludedPaths []string              `json:"excluded_paths"`
	RuleConfigs   map[string]RuleConfig `json:"rule_configs"`
}

func (c RuntimeConfig) match(path string) bool {
	for _, pattern := range c.ExcludedPaths {
		if matched, err := filepath.Match(pattern, path); matched || err != nil {
			return false
		}
	}
	for _, pattern := range c.IncludedPaths {
		if matched, err := filepath.Match(pattern, path); matched && err == nil {
			return true
		}
	}
	return false
}
