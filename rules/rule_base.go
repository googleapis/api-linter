package rules

import "github.com/jgeewax/api-linter/lint"

// ruleBase implements lint.Rule.
type ruleBase struct {
	RuleInfo RuleInfo
	l        Linter
}

func (r ruleBase) Name() string {
	return r.RuleInfo.Name
}

func (r ruleBase) Description() string {
	return r.RuleInfo.Description
}

func (r ruleBase) URL() string {
	return r.RuleInfo.URL
}

func (r ruleBase) FileTypes() []lint.FileType {
	return r.RuleInfo.FileTypes
}

func (r ruleBase) Category() lint.Category {
	return r.RuleInfo.Category
}

func (r ruleBase) Lint(req lint.Request) (lint.Response, error) {
	return r.l.Lint(req, r)
}

// RuleInfo stores information of a rule.
type RuleInfo struct {
	Name        string          // rule name in the set.
	Description string          // a short description of this rule.
	URL         string          // a link to a document for more details.
	FileTypes   []lint.FileType // types of files that this rule targets to.
	Category    lint.Category   // category of problems this rule produces.
}
