package lint

// ruleInfo stores information of a rule.
type RuleInfo struct {
	Name        string     // rule name in the set.
	Description string     // a short description of this rule.
	Url         string     // a link to a document for more details.
	FileTypes   []FileType // types of files that this rule targets to.
	Category    Category   // category of problems this rule produces.

	noPositional struct{} // Prevent positional composite literal instantiation
}
