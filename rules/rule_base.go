package rules

import "github.com/jgeewax/api-linter/lint"

// linter is the interface that wraps lint function.
type linter interface {
	Lint(rule lint.Rule, req lint.Request) (lint.Response, error)
}

// ruleBase implements lint.Rule.
type ruleBase struct {
	metadata metadata
	l        linter
}

func (r ruleBase) ID() lint.RuleID {
	return lint.RuleID{
		Set:  r.metadata.Set,
		Name: r.metadata.Name,
	}
}

func (r ruleBase) Description() string {
	return r.metadata.Description
}

func (r ruleBase) URL() string {
	return r.metadata.URL
}

func (r ruleBase) FileTypes() []lint.FileType {
	return r.metadata.FileTypes
}

func (r ruleBase) Category() lint.Category {
	return r.metadata.Category
}

func (r ruleBase) Lint(req lint.Request) (lint.Response, error) {
	return r.l.Lint(r, req)
}

// metadata stores metadata information of a lint.Rule.
type metadata struct {
	Set         string          // rule set name.
	Name        string          // rule name in the set.
	Description string          // a short description of this rule.
	URL         string          // a link to a document for more details.
	FileTypes   []lint.FileType // types of files that this rule targets to.
	Category    lint.Category   // category of problems this rule produces.
}
