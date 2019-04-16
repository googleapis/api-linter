package lint

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule interface {
	// Info returns metadata about a rule.
	Info() RuleInfo
	// Lint performs the linting process.
	Lint(Request) (Response, error)
}

// FileType defines a file type that a rule is targeting to.
type FileType string

const (
	// ProtoFile indicates that the targeted file is a protobuf-definition file.
	ProtoFile FileType = "proto-file"
)

// Category defines the category of the findings produced by a rule.
type Category string

const (
	// CategoryError indicates that in the API, something will cause errors.
	CategoryError Category = "API-Linter-Error"
	// CategorySuggestion indicates that in the API, something can be do better.
	CategorySuggestion Category = "API-Linter-Suggestion"
)
