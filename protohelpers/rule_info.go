package protohelpers

import "github.com/jgeewax/api-linter/lint"

// ruleInfo stores information of a rule.
type RuleInfo struct {
	name        string          // rule name in the set.
	description string          // a short description of this rule.
	url         string          // a link to a document for more details.
	fileTypes   []lint.FileType // types of files that this rule targets to.
	category    lint.Category   // category of problems this rule produces.
}

func NewRuleInfo(name, description, url string, fileTypes []lint.FileType, category lint.Category) RuleInfo {
	return RuleInfo{
		name:        name,
		description: description,
		url:         url,
		fileTypes:   fileTypes,
		category:    category,
	}
}

func (r RuleInfo) Name() string {
	return r.name
}

func (r RuleInfo) Description() string {
	return r.description
}

func (r RuleInfo) URL() string {
	return r.url
}

func (r RuleInfo) FileTypes() []lint.FileType {
	return r.fileTypes
}

func (r RuleInfo) Category() lint.Category {
	return r.category
}
