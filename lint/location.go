package lint

// Location describes a location in a source code file.
type Location struct {
	Start, End Position
}

// IsValid checks if the location struct is constructed properly and
// returns true if it's valid.
func (l Location) IsValid() bool {
	return l.Start.IsValid() && l.End.IsValid()
}

// Position describes a zero-based position in a source code file.
// Typically you will want to add 1 to each before displaying to a user.
type Position struct {
	Line, Column int
}

// IsValid checks if the position struct is constructed properly and
// returns true if it's valid.
func (p Position) IsValid() bool {
	return p.Line >= 0 && p.Column >= 0
}
