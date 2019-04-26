package lint

import (
	"errors"
	"strings"
)

// Runtime stores a set of rules and configs.
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

// AddRules adds rules.
//
// Note: it will check name conflict.
func (r *Runtime) AddRules(rules ...Rule) error {
	for _, rl := range rules {
		if err := r.rules.Register(rl); err != nil {
			return err
		}
	}
	return nil
}

// Run executes rules on the request.
//
// It uses the proto file path to determine which rules will
// be applied to the request, according to the list of runtime
// configs.
func (r *Runtime) Run(req Request) (Response, error) {
	finalResp := Response{}
	var errMessages []string

	for name, rl := range r.rules {
		config := getDefaultRuleConfig()

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
