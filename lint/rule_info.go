package lint

// RuleInfo defines meta-information for a Rule.
type RuleInfo interface {
	Name() string          // returns the rule name.
	Description() string   // returns a short description of this rule.
	URL() string           // returns a link to a document for more details.
	FileTypes() []FileType // returns the list of FileType that this rule is targeting to.
	Category() Category    // returns the Category this rule will produce.
}

// ruleInfo stores information of a rule.
type ruleInfo struct {
	name        string     // rule name in the set.
	description string     // a short description of this rule.
	url         string     // a link to a document for more details.
	fileTypes   []FileType // types of files that this rule targets to.
	category    Category   // category of problems this rule produces.
}

// NewRuleInfo creates and returns a RuleInfo from the provided information.
func NewRuleInfo(name, description, url string, fileTypes []FileType, category Category) RuleInfo {
	return ruleInfo{
		name:        name,
		description: description,
		url:         url,
		fileTypes:   fileTypes,
		category:    category,
	}
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

func (r ruleInfo) FileTypes() []FileType {
	return r.fileTypes
}

func (r ruleInfo) Category() Category {
	return r.category
}
