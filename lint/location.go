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

// Position describes a one-based position in a source code file.
type Position struct {
	Line, Column int
}

// IsValid checks if the position struct is constructed properly and
// returns true if it's valid.
func (p Position) IsValid() bool {
	return p.Line > 0 && p.Column > 0
}

// `FileStartLocation` returns a `Location` representing the starting point of a file
func FileStartLocation() Location {
	return Location{
		Start: Position{1, 1},
		End:   Position{1, 1},
	}
}
