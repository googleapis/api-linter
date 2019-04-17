package lint

import (
	"testing"
)

func TestLocation_IsValid(t *testing.T) {
	tests := []struct {
		l     Location
		valid bool
	}{
		{Location{Position{1, 2}, Position{1, 2}}, true},
		{invalidLocation(), false},
		{Location{Position{-1, -1}, Position{0, 0}}, false}, // invalid Start
		{Location{Position{0, 0}, Position{-1, -1}}, false}, // invalid End
		{Location{Position{1, 1}, Position{0, 0}}, false},   // End is before Start
		{Location{Position{1, 1}, Position{1, 0}}, false},   // End is before Start
	}

	for _, test := range tests {
		v := test.l.IsValid()

		if v != test.valid {
			t.Errorf("%+v.IsValid()=%t; want %t", test.l, v, test.valid)
		}
	}
}
