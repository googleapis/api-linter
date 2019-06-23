package protos

import "testing"

func TestFindByName(t *testing.T) {
	tests := []struct {
		name      string
		mustFound bool
	}{
		{"google/type/date.proto", true},
		{"google/testing_not_found.proto", false},
	}

	for _, test := range tests {
		f, err := FindByName(test.name)
		if test.mustFound && err != nil {
			t.Errorf("FindByName(%q) returns error: %v", test.name, err)
		}
		if !test.mustFound && err == nil {
			t.Errorf("FindByName(%q): expect not found error, but got nil", test.name)
		}
		if err == nil && f.GetName() != test.name {
			t.Errorf("FindByName(%q) returns proto file with name %q", test.name, f.GetName())
		}
	}
}

func TestFindAllByNames(t *testing.T) {
	tests := []struct {
		names     []string
		mustFound bool
	}{
		{[]string{"google/type/date.proto"}, true},
		{[]string{"google/testing_not_found.proto"}, false},
	}

	for _, test := range tests {
		files, err := FindAllByNames(test.names...)
		if test.mustFound && err != nil {
			t.Errorf("FindMoreByNames(%v) returns error: %v", test.names, err)
		}
		if !test.mustFound && err == nil {
			t.Errorf("FindMoreByNames(%v): expect not found error, but got nil", test.names)
		}
		if err == nil && len(files) != len(test.names) {
			t.Errorf("FindMoreByNames(%v) returns %d files, but want %d", test.names, len(files), len(test.names))
		}
	}
}
