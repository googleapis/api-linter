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
func NewRuntime(c ...RuntimeConfig) *Runtime {
	t := &Runtime{
		rules: make(Rules),
	}
	t.AddConfigs(c...)
	return t
}

// AddConfigs adds a runtime config.
func (r *Runtime) AddConfigs(c ...RuntimeConfig) {
	r.configs = append(r.configs, c...)
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
	resp := Response{
		FilePath: req.ProtoFile().Path(),
	}
	var errMessages []string

	for name, rl := range r.rules {
		config := getDefaultRuleConfig()

		if c, err := r.configs.getRuleConfig(req.ProtoFile().Path(), name); err == nil {
			config = config.withOverride(c)
		} else {
			errMessages = append(errMessages, err.Error())
			continue
		}

		if config.Status == Enabled && !req.DescriptorSource().isRuleDisabledInFile(rl.Info().Name) {
			if problems, err := rl.Lint(req); err == nil {
				for _, p := range problems {
					if !req.DescriptorSource().isRuleDisabled(rl.Info().Name, p.Descriptor) {
						p.RuleID = rl.Info().Name
						p.Category = config.Category
						resp.Problems = append(resp.Problems, p)
					}
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

	return resp, err
}
