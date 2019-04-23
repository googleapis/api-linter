package lint

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Configs stores a list of Config and supports config lookup
// for a given path.
type Configs []Config

// Search returns the first found config that matches the path
// or an error if not found.
func (c Configs) Search(path string) (*Config, error) {
	for _, cfg := range c {
		if cfg.match(path) {
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("no config matches path %q", path)
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

// RuleConfig stores the status and category of a rule,
// which can be applied to a rule during runtime.
type RuleConfig struct {
	Status   Status   `json:"status"`
	Category Category `json:"category"`
}

// Config stores rule runtime configurations and file spec that
// a path must match any of the included paths but none of
// the excluded paths.
type Config struct {
	IncludedPaths []string              `json:"included_paths"`
	ExcludedPaths []string              `json:"excluded_paths"`
	RuleConfigs   map[string]RuleConfig `json:"rule_configs"`
}

func (c Config) match(path string) bool {
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
