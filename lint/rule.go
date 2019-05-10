package lint

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule interface {
	// Info returns metadata about a rule.
	Info() RuleInfo
	// Lint performs the linting process.
	Lint(Request) ([]Problem, error)
}
