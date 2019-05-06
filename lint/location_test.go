package lint

import (
	"testing"
)

func TestLocation_IsValid(t *testing.T) {
	tests := []struct {
		l     *Location
		valid bool
	}{
		{&Location{&Position{1, 2}, &Position{1, 2}}, true},
		{&Location{&Position{-1, -1}, &Position{0, 0}}, false}, // invalid: start position is invalid
		{&Location{&Position{0, 0}, &Position{-1, -2}}, false}, // invalid: end position is invalid
		{&Location{&Position{1, 1}, &Position{0, 0}}, false},   // invalid: end line is before start line
		{&Location{&Position{1, 1}, &Position{1, 0}}, false},   // invalid: end column is before start column in the same line
		{&Location{&Position{1, 1}, &Position{1, 0}}, false},   // invalid: start position is not created properly
		{nil, false}, // invalid: nil
	}

	for _, test := range tests {
		v := test.l.IsValid()

		if v != test.valid {
			t.Errorf("%+v.IsValid()=%t; want %t", test.l, v, test.valid)
		}
	}
}

func TestPosition_IsValid(t *testing.T) {
	tests := []struct {
		p     *Position
		valid bool
	}{
		{&Position{0, 1}, true},
		{&Position{-1, 0}, false},
		{&Position{0, -1}, false},
		{nil, false},
	}

	for _, test := range tests {
		if got, want := test.p.IsValid(), test.valid; got != want {
			t.Errorf("Position.IsValid() returns %v, but want %v", got, want)
		}
	}
}
