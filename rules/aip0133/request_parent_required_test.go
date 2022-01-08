package aip0133

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestRequestParentFieldRequired(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		name        string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "CreateBookRequest", "parent", nil},
		{"MissingParent", "CreateBookRequest", "id", testutils.Problems{{Message: "parent"}}},
		{"IrrelevantMessage", "AddBookRequest", "id", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					string {{.FieldName}} = 1;
				}
			`, test)

			problems := requestParentRequired.Lint(f)
			message := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}

	// Test the "top-level exception", which is more involved
	// than the other tests and therefore handled separately.
	t.Run("ValidTopLevel", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message CreateBookRequest {
				Book book = 2;
			}
			message Book {
				option (google.api.resource) = {
					pattern: "books/{book}"
				};
				string name = 1;
			}
		`)
		problems := requestParentRequired.Lint(f)
		if diff := (testutils.Problems{}).Diff(problems); diff != "" {
			t.Errorf(diff)
		}
	})
}
