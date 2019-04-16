package rules

import (
	"github.com/jgeewax/api-linter/lint"
	"testing"
)

func TestFieldNamesRule_ConformingFieldNames(t *testing.T) {
	pd := protoDescriptorProtoFromSource([]byte(`syntax = "proto2";

package google.apis.tools.analyzer.testprotos;

message Outer {
  message Middle {
      optional string middle_field_name = 1;
  }

  optional string outer_field_name = 2;

  enum NestedEnum {
    FOO = 1;
  }
  oneof outer_oneof_field {
    string outer_oneof_field_name = 3;
  }

  extensions 100 to 199;
}`))

	rules := lint.NewRules()
}
