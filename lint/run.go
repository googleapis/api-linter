package lint

import (
	"errors"
	"strings"
)

// Run invokes all rules on the request.
func Run(rules Rules, request Request) (Response, error) {
	return run(rules.AllRules(), request)
}

// RunWithConfig invokes rules filtered by the `RulesConfig` on the request.
func RunWithConfig(rules Rules, request Request, cfg RulesConfig) (Response, error) {
	return run(rules.FindRulesByConfig(cfg), request)
}

func run(rules []Rule, req Request) (Response, error) {
	finalResp := Response{}
	errMsgs := []string{}
	for _, r := range rules {
		if resp, err := r.Lint(req); err == nil {
			finalResp.Merge(resp)
		} else {
			errMsgs = append(errMsgs, err.Error())
		}
	}

	if len(errMsgs) != 0 {
		err := errors.New(strings.Join(errMsgs, "; "))
		return finalResp, err
	}

	return finalResp, nil
}
