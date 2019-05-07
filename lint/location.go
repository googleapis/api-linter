package lint

import "fmt"

// Location describes a location in a source code file.
type Location struct {
	Start Position `json:"start_position" yaml:"start_position"`
	End   Position `json:"end_position" yaml:"end_position"`
}

// IsValid checks if the location is constructed properly and
// returns true if so.
func (l *Location) IsValid() bool {
	return l != nil &&
		l.Start.IsValid() &&
		l.End.IsValid() &&
		(l.End.Line > l.Start.Line ||
			l.End.Line == l.Start.Line && l.End.Column >= l.Start.Column)
}

// String returns the string representation.
func (l *Location) String() string {
	return fmt.Sprintf("{start: %s, end: %s}", l.Start, l.End)
}

// Position describes a zero-based position in a source code file.
// Typically you will want to add 1 to each before displaying to a user.
type Position struct {
	Line   int `json:"line_number" yaml:"line_number"`
	Column int `json:"column_number" yaml:"column_number"`
}

// IsValid checks if the position is constructed properly and
// returns true if so.
func (p Position) IsValid() bool {
	return p.Line >= 0 && p.Column >= 0
}

// String returns the string representation.
func (p Position) String() string {
	return fmt.Sprintf("{line: %d, column: %d}", p.Line, p.Column)
}
