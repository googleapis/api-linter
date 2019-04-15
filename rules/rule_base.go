package rules

import "github.com/jgeewax/api-linter/lint"

// ruleBase implements lint.Rule.
type ruleBase struct {
	ruleInfo
	l linter
}

func (r ruleBase) Lint(req lint.Request) (lint.Response, error) {
	return r.l.Lint(req, r)
}

// ruleInfo stores information of a rule.
type ruleInfo struct {
	name        string          // rule name in the set.
	description string          // a short description of this rule.
	url         string          // a link to a document for more details.
	fileTypes   []lint.FileType // types of files that this rule targets to.
	category    lint.Category   // category of problems this rule produces.
}

func (r ruleInfo) Name() string {
	return r.name
}

func (r ruleInfo) Description() string {
	return r.description
}

func (r ruleInfo) URL() string {
	return r.url
}

func (r ruleInfo) FileTypes() []lint.FileType {
	return r.fileTypes
}

func (r ruleInfo) Category() lint.Category {
	return r.category
}
