package lint

import (
	"testing"
)

func TestLocation_IsValid(t *testing.T) {
	tests := []struct {
		l     *Location
		valid bool
	}{
		{NewLocation(NewPosition(1, 2), NewPosition(1, 2)), true},
		{NewLocation(NewPosition(-1, -1), NewPosition(0, 0)), false}, // invalid: start position is invalid
		{NewLocation(NewPosition(0, 0), NewPosition(-1, -1)), false}, // invalid: end position is invalid
		{NewLocation(NewPosition(1, 1), NewPosition(0, 0)), false},   // invalid: end line is before start line
		{NewLocation(NewPosition(1, 1), NewPosition(1, 0)), false},   // invalid: end column is before start column in the same line
		{NewLocation(&Position{}, NewPosition(1, 0)), false},         // invalid: start position is not created properly
		{NewLocation(NewPosition(1, 0), &Position{}), false},         // invalid: end position is not created properly
		{&Location{}, false}, // invalid: location is not created properly
		{nil, false},         // invalid: nil
	}

	for _, test := range tests {
		v := test.l.IsValid()

		if v != test.valid {
			t.Errorf("%+v.IsValid()=%t; want %t", test.l, v, test.valid)
		}
	}
}

func TestLocation_Start(t *testing.T) {
	pos := NewPosition(0, 0)
	tests := []struct {
		loc *Location
		pos *Position
	}{
		{NewLocation(pos, nil), pos},
		{NewLocation(nil, pos), nil},
		{nil, nil},
	}

	for _, test := range tests {
		if got, want := test.loc.Start(), test.pos; got != want {
			t.Errorf("Location.Start() returns %v, but want %v", got, want)
		}
	}
}

func TestLocation_End(t *testing.T) {
	pos := NewPosition(0, 0)
	tests := []struct {
		loc *Location
		pos *Position
	}{
		{NewLocation(nil, pos), pos},
		{NewLocation(pos, nil), nil},
		{nil, nil},
	}

	for _, test := range tests {
		if got, want := test.loc.End(), test.pos; got != want {
			t.Errorf("Location.End() returns %v, but want %v", got, want)
		}
	}
}

func TestPosition_Line(t *testing.T) {
	tests := []struct {
		p *Position
		l int
	}{
		{NewPosition(1, 0), 1},
		{&Position{}, -1},
		{nil, -1},
	}

	for _, test := range tests {
		if got, want := test.p.Line(), test.l; got != want {
			t.Errorf("Position.Line() returns %d, but want %d", got, want)
		}
	}
}

func TestPosition_Column(t *testing.T) {
	tests := []struct {
		p *Position
		c int
	}{
		{NewPosition(0, 1), 1},
		{&Position{}, -1},
		{nil, -1},
	}

	for _, test := range tests {
		if got, want := test.p.Column(), test.c; got != want {
			t.Errorf("Position.Column() returns %d, but want %d", got, want)
		}
	}
}

func TestPosition_IsValid(t *testing.T) {
	tests := []struct {
		p     *Position
		valid bool
	}{
		{NewPosition(0, 1), true},
		{NewPosition(-1, 0), false},
		{NewPosition(0, -1), false},
		{&Position{}, false},
		{nil, false},
	}

	for _, test := range tests {
		if got, want := test.p.IsValid(), test.valid; got != want {
			t.Errorf("Position.IsValid() returns %v, but want %v", got, want)
		}
	}
}
