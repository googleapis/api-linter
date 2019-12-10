package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestResourceFieldRequired(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		name              string
		MessageName       string
		ResourceName      string
		ResourceFieldName string
		problems          testutils.Problems
	}{
		{"Valid", "UpdateBookRequest", "Book", "book", nil},
		{"ValidTwoWords", "UpdateBigBookRequest", "BigBook", "big_book", nil},
		{"InvalidMismatch", "UpdateBookRequest", "Foo", "foo", testutils.Problems{{Message: "has no \"Book\""}}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message {{.MessageName}} {
					{{.ResourceName}} {{.ResourceFieldName}} = 1;
				}
				message {{.ResourceName}} {}
			`, test)
			message := file.GetMessageTypes()[0]
			problems := requestResourceFieldRequired.Lint(file)
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
