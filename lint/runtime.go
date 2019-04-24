package lint

import (
	"errors"
	"fmt"
	"strings"
)

// Runtime stores a set of rules.
type Runtime struct {
	rules  Rules
	config RuleConfig
}

// NewRuntime creates a new Runtime.
func NewRuntime(config RuleConfig) *Runtime {
	return &Runtime{
		rules:  make(Rules),
		config: config,
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

// Run executes rules on the request when a config is found for the file path
// of the request.
//
// If the found config contains rule configs for some rules, the status and
// category of the affected rules will be updated accordingly. In other words,
// rule configs can be used to turn on/off certain rules and change the category
// of the returned problems.
func (r *Runtime) Run(req Request, configs RuntimeConfigs) (Response, error) {
	cfg, err := configs.Search(req.ProtoFile().Path())

	if err != nil {
		return Response{}, err
	}

	finalResp := Response{}
	var errMessages []string
	for name, rl := range r.rules {
		config := r.config
		for prefix, c := range cfg.RuleConfigs {
			if name.HasPrefix(prefix) {
				config = config.WithOverride(c)
				break
			}
		}

		if config.Status == Enabled {
			if resp, err := rl.Lint(req); err == nil {
				for _, p := range resp.Problems {
					p.category = rl.Info().Category

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

	if len(errMessages) != 0 {
		err = errors.New(strings.Join(errMessages, "; "))
	}

	return finalResp, err
}
