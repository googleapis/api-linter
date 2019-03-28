package lint

// Location describes a location in a source code file.
type Location struct {
	Start, End Position
}

// Position describes a one-based position in a source code file.
type Position struct {
	Line, Column int
}

// StartLocation marks the very start position in a file.
var StartLocation = Location{
	Start: Position{1, 1},
	End:   Position{1, 1},
}
