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
func (r *Runtime) AddRules(prefix string, rules ...Rule) error {
	for _, rl := range rules {
		if _, found := r.rules[rl.Info().Name.WithPrefix(prefix)]; found {
			return fmt.Errorf("duplicate repository entry with name %q", rl.Info().Name.WithPrefix(prefix))
		}
		r.rules[rl.Info().Name.WithPrefix(prefix)] = rl
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
		config, err := r.configs.getRuleConfig(req.ProtoFile().Path(), name)

		if err != nil {
			errMessages = append(errMessages, err.Error())
			continue
		}

		config = defaultRuleConfig.withOverride(config)

		if config.Status == Enabled {
			if resp, err := rl.Lint(req); err == nil {
				for _, p := range resp.Problems {

					if config.Category != "" {
						p.category = config.Category
					}

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
