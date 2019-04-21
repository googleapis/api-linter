package lint

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// NameSeparator denotes the separator in the rule name.
	NameSeparator string = "::"
)

// Repository stores a set of rules.
type Repository struct {
	ruleMap map[string]Rule
}

// NewRepository creates a new Repository.
func NewRepository() *Repository {
	return &Repository{
		ruleMap: make(map[string]Rule),
	}
}

// AddRule adds rules, of which the name will be added a prefix to
// reduce conflict, and the status and category will be changed
// by the given rule config.
func (r *Repository) AddRule(prefix string, cfg RuleConfig, rule ...Rule) error {
	for _, rl := range rule {
		rl.Info().Name = prefix + NameSeparator + rl.Info().Name
		if cfg.Status != "" {
			rl.Info().Status = cfg.Status
		}
		if cfg.Category != "" {
			rl.Info().Category = cfg.Category
		}

		if _, found := r.ruleMap[rl.Info().Name]; found {
			return fmt.Errorf("duplicate rule name `%s`", rl.Info().Name)
		}
		r.ruleMap[rl.Info().Name] = rl
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
	finalResp := Response{}
	errMessages := []string{}
	for name, rl := range r.ruleMap {
		ruleCfg := RuleConfig{
			Status:   rl.Info().Status,
			Category: rl.Info().Category,
		}
		for prefix, c := range ruleCfgMap {
			if strings.HasPrefix(name, prefix) {
				ruleCfg = c
				break
			}
		}
		if ruleCfg.Status == Enabled {
			if resp, err := rl.Lint(req); err == nil {
				for _, p := range resp.Problems {
					p.Category = ruleCfg.Category
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
