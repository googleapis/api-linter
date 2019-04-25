package lint

import (
	"errors"
	"fmt"
	"strings"
)

// Runtime stores a set of rules.
type Runtime struct {
	rules   Rules
	configs RuntimeConfigs
}

// NewRuntime creates a new Runtime.
func NewRuntime(configs RuntimeConfigs) *Runtime {
	return &Runtime{
		rules:   make(Rules),
		configs: configs,
	}
}

// AddRules adds rules, of which the name will be added a prefix to reduce collisions
func (r *Runtime) AddRules(rules ...Rule) error {
	for _, rl := range rules {
		if _, found := r.rules[rl.Info().Name]; found {
			return fmt.Errorf("duplicate repository entry with name %q", rl.Info().Name)
		}
		r.rules[rl.Info().Name] = rl
	}
	return nil
}

// Run executes rules on the request when a config is found for the file path of the request.
//
// If the found config contains rule configs for some rules, the status and
// category of the affected rules will be updated accordingly. In other words,
// rule configs can be used to turn on/off certain rules and change the category
// of the returned problems.
func (r *Runtime) Run(req Request) (Response, error) {
	finalResp := Response{}
	var errMessages []string

	for name, rl := range r.rules {
		config := defaultRuleConfig

		if c, err := r.configs.getRuleConfig(req.ProtoFile().Path(), name); err == nil {
			config = config.withOverride(c)
		} else {
			errMessages = append(errMessages, err.Error())
			continue
		}

		if config.Status == Enabled {
			if resp, err := rl.Lint(req); err == nil {
				for _, p := range resp.Problems {
					p.category = config.Category
					finalResp.Problems = append(finalResp.Problems, p)
				}
			} else {
				errMessages = append(errMessages, err.Error())
			}
		}
	}

	var err error
	if len(errMessages) != 0 {
		err = errors.New(strings.Join(errMessages, "; "))
	}

	return finalResp, err
}
