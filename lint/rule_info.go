package lint

import (
	"regexp"
	"strings"
)

// RuleInfo stores information of a rule.
type RuleInfo struct {
	Name         RuleName      // rule name in the set.
	Description  string        // a short description of this rule.
	URI          string        // a link to a document for more details.
	RequestTypes []RequestType // types of requests that this rule should receive.

	noPositional struct{} // Prevent positional composite literal instantiation
}

// RequestType defines a request type that can be passed into rules.
type RequestType string

const (
	// ProtoRequest indicates that the targeted request contains a protobuf-definition file.
	ProtoRequest RequestType = "proto-request"
)

// Status defines whether a rule is enabled, disabled or deprecated.
type Status string

const (
	// Enabled indicates that a rule should be enabled.
	Enabled Status = "enabled"
	// Disabled indicates that a rule should be disabled.
	Disabled Status = "disabled"
	// Deprecated indicates that a rule should be deprecated.
	Deprecated Status = "Deprecated"
)

// RuleName is an identifier for a rule. Allowed characters include a-z, A-Z, 0-9, -, _. The
// namespace separator :: is allowed between RuleName segments (for example, my_namespace::my_rule).
type RuleName string

const nameSeparator string = "::"

var ruleNameValidator = regexp.MustCompile("^([a-zA-Z0-9-_]+(::[a-zA-Z0-9-_]+)?)+$")

// NewRuleName creates a RuleName from segments. It will join the segments with the "::" separator.
func NewRuleName(segments ...string) RuleName {
	return RuleName(strings.Join(segments, nameSeparator))
}

// IsValid checks if a RuleName is syntactically valid.
func (r RuleName) IsValid() bool {
	return r != "" && ruleNameValidator.Match([]byte(r))
}

func (r RuleName) parent() RuleName {
	lastSeparator := strings.LastIndex(string(r), nameSeparator)

	if lastSeparator == -1 {
		return ""
	}

	return r[:lastSeparator]
}

// HasPrefix returns true if r contains prefix as a namespace. prefix parameters can be "::" delimited
// or specified as independent parameters.
// For example:
//
// r := NewRuleName("foo", "bar", "baz")   // string(r) == "foo::bar::baz"
//
// r.HasPrefix("foo::bar")          == true
// r.HasPrefix("foo", "bar")        == true
// r.HasPrefix("foo", "bar", "baz") == true   // matches the entire string
// r.HasPrefix("foo", "ba")         == false  // prefix must end on a delimiter
func (r RuleName) HasPrefix(prefix ...string) bool {
	prefixSegments := make([]string, 0, len(prefix))

	for _, prefixSegment := range prefix {
		prefixSegments = append(prefixSegments, strings.Split(prefixSegment, "::")...)
	}

	prefixStr := strings.Join(prefixSegments, nameSeparator)

	return string(r) == prefixStr || strings.HasPrefix(string(r), prefixStr+nameSeparator)
}
