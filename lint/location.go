package lint

import "fmt"

// Location describes a location in a source code file.
//
// Note: positions are one-based.
type Location struct {
	Start Position `json:"start_position" yaml:"start_position"`
	End   Position `json:"end_position" yaml:"end_position"`
}

// IsValid checks if the location is constructed properly and
// returns true if so.
func (l Location) IsValid() bool {
	return l.Start.IsValid() &&
		l.End.IsValid() &&
		(l.End.Line > l.Start.Line ||
			l.End.Line == l.Start.Line && l.End.Column >= l.Start.Column)
}

// String returns the string representation.
func (l Location) String() string {
	return fmt.Sprintf("{start: %s, end: %s}", l.Start, l.End)
}

// Position describes a one-based position in a source code file.
type Position struct {
	Line   int `json:"line_number" yaml:"line_number"`
	Column int `json:"column_number" yaml:"column_number"`
}

// IsValid checks if the position is constructed properly and
// returns true if so.
func (p Position) IsValid() bool {
	return p.Line >= 1 && p.Column >= 1
}

// String returns the string representation.
func (p Position) String() string {
	return fmt.Sprintf("{line: %d, column: %d}", p.Line, p.Column)
}
