package lint

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule interface {
	// Name returns an unique name for this rule.
	Name() string
	// Description returns a short summary about this rule.
	Description() string
	// URL returns a link at which the full documentation
	// about this rule can be accessed.
	URL() string
	// FileTypes returns the types of files that this rule is targeting to.
	// E.g., `ProtoFile` for protobuf files.
	FileTypes() []FileType
	// Category returns the category of findings produced
	// by this rule, e.g. Problem, Suggestion, etc.
	Category() Category
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
