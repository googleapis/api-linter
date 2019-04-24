package lint

import (
	"errors"
	"fmt"
	"strings"
)

// Repository stores a set of rules.
type Repository struct {
	ruleMap map[RuleName]ruleEntry
}

type ruleEntry struct {
	name            RuleName
	rule            Rule
	defaultCategory Category
	status          Status
}

// NewRepository creates a new Repository.
func NewRepository() *Repository {
	return &Repository{
		ruleMap: make(map[RuleName]ruleEntry),
	}
}

// AddRule adds rules, of which the name will be added a prefix to
// reduce conflict, and the status and category will be changed
// by the given rule config.
func (r *Repository) AddRule(prefix string, cfg RuleConfig, rule ...Rule) error {
	for _, rl := range rule {
		e := ruleEntry{
			name:            rl.Info().Name.WithPrefix(prefix),
			rule:            rl,
			defaultCategory: rl.Info().Category,
			status:          Enabled,
		}
		if cfg.Status != "" {
			e.status = cfg.Status
		}
		if cfg.Category != "" {
			e.defaultCategory = cfg.Category
		}

		if _, found := r.ruleMap[e.name]; found {
			return fmt.Errorf("duplicate repository entry with name %q", e.name)
		}
		r.ruleMap[e.name] = e
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
func (r *Repository) Run(req Request, configs Configs) (Response, error) {
	cfg, err := configs.Search(req.ProtoFile().Path())
	if err != nil {
		return Response{}, err
	}
	return r.run(req, cfg.RuleConfigs)
}

func (r *Repository) run(req Request, ruleCfgMap map[string]RuleConfig) (Response, error) {
	disabledRules := newDisabledRuleFinder(req.ProtoFile(), req.DescriptorSource())
	finalResp := Response{}
	errMessages := []string{}
	for name, rl := range r.ruleMap {
		ruleCfg := RuleConfig{
			Status:   rl.status,
			Category: rl.defaultCategory,
		}
		for prefix, c := range ruleCfgMap {
			if name.HasPrefix(prefix) {
				if c.Status != "" {
					ruleCfg.Status = c.Status
				}
				if c.Category != "" {
					ruleCfg.Category = c.Category
				}
				break
			}
		}
		if ruleCfg.Status == Enabled {
			if resp, err := rl.rule.Lint(req); err == nil {
				for _, p := range resp.Problems {
					ruleDisabled := false
					if p.Location != nil {
						ruleDisabled = disabledRules.isRuleDisabledAtLocation(name, p.Location)
					} else {
						ruleDisabled = disabledRules.isRuleDisabledAtDescriptor(name, p.Descriptor)
					}

					if !ruleDisabled {
						p.category = ruleCfg.Category
						finalResp.Problems = append(finalResp.Problems, p)
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

	return finalResp, err
}
