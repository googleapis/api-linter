package rules

import "github.com/jgeewax/api-linter/lint"

// ruleBase implements lint.Rule.
type ruleBase struct {
	ruleInfo ruleInfo
	l        linter
}

func (r ruleBase) Name() string {
	return r.ruleInfo.Name
}

func (r ruleBase) Description() string {
	return r.ruleInfo.Description
}

func (r ruleBase) URL() string {
	return r.ruleInfo.URL
}

func (r ruleBase) FileTypes() []lint.FileType {
	return r.ruleInfo.FileTypes
}

func (r ruleBase) Category() lint.Category {
	return r.ruleInfo.Category
}

func (r ruleBase) Lint(req lint.Request) (lint.Response, error) {
	return r.l.Lint(req, r)
}

// ruleInfo stores information of a rule.
type ruleInfo struct {
	Name        string          // rule name in the set.
	Description string          // a short description of this rule.
	URL         string          // a link to a document for more details.
	FileTypes   []lint.FileType // types of files that this rule targets to.
	Category    lint.Category   // category of problems this rule produces.
}
