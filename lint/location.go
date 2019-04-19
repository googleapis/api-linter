package lint

import "fmt"

// Location describes a location in a source code file.
// Use `NewLocation` to create a new Location -- Location{}
// will be invalid.
type Location struct {
	start, end *Position
}

// NewLocation creates a new Location from the start and end positions.
func NewLocation(start, end *Position) *Location {
	return &Location{
		start: start,
		end:   end,
	}
}

// Start returns the start position.
func (l *Location) Start() *Position {
	if l == nil {
		return nil
	}
	return l.start
}

// End returns the end position.
func (l *Location) End() *Position {
	if l == nil {
		return nil
	}
	return l.end
}

// IsValid checks if the location is constructed properly and
// returns true if so.
func (l *Location) IsValid() bool {
	return l != nil &&
		l.start.IsValid() &&
		l.end.IsValid() &&
		(l.end.Line() > l.start.Line() ||
			l.end.Line() == l.start.Line() && l.end.Column() >= l.start.Column())
}

// String returns the string representation.
func (l *Location) String() string {
	return fmt.Sprintf("{start: %s, end: %s}", l.Start(), l.End())
}

// Position describes a zero-based position in a source code file.
// Typically you will want to add 1 to each before displaying to a user.
// Use `NewPosition` to create a new Position -- `Position{}` will be
// invalid.
type Position struct {
	line, column int

	// Since line and column are both zero-based, we need to check if
	// the position is constructed properly using the `NewPosition` function.
	filled bool
}

// NewPosition returns a new Position from the line and column numbers.
func NewPosition(line, column int) *Position {
	return &Position{
		line:   line,
		column: column,
		filled: true,
	}
}

func (p *Position) isFilled() bool {
	return p != nil && p.filled
}

// Line returns the zero-based line number.
// It returns a negative number if the position is not valid.
// Always check `IsVaild` before using this.
func (p *Position) Line() int {
	if !p.isFilled() {
		return -1
	}
	return p.line
}

// Column returns the zero-based column number.
// It returns a negative number if the position is not valid.
// Always check `IsValid` before using this.
func (p *Position) Column() int {
	if !p.isFilled() {
		return -1
	}
	return p.column
}

// IsValid checks if the position is constructed properly and
// returns true if so.
func (p *Position) IsValid() bool {
	return p.isFilled() && p.Line() >= 0 && p.Column() >= 0
}

// String returns the string representation.
func (p *Position) String() string {
	return fmt.Sprintf("{line: %d, column: %d}", p.Line(), p.Column())
}
