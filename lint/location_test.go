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
