// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0148

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestIpAddressFormat(t *testing.T) {
	for _, test := range []struct {
		name, Annotation, FieldName, Type string
		problems                          testutils.Problems
	}{
		{
			name:       "ValidIpAddressFormat",
			FieldName:  "ip_address",
			Annotation: "[(google.api.field_info).format = IPV4_OR_IPV6]",
			Type:       "string",
		},
		{
			name:       "ValidIpAddressSuffixFormat",
			FieldName:  "incoming_ip_address",
			Annotation: "[(google.api.field_info).format = IPV4_OR_IPV6]",
			Type:       "string",
		},
		{
			name:      "SkipNonIpAddress",
			FieldName: "other",
			Type:      "string",
		},
		{
			name:      "SkipNonIpAddressSuffix",
			FieldName: "zip_address",
			Type:      "string",
		},
		{
			name:      "SkipNonStringIpAddress",
			FieldName: "ip_address",
			Type:      "Foo",
		},
		{
			name:      "InvalidMissingFormat",
			FieldName: "ip_address",
			Type:      "string",
			problems:  testutils.Problems{{Message: "must specify one of"}},
		},
		{
			name:       "InvalidWrongFormat",
			FieldName:  "ip_address",
			Annotation: "[(google.api.field_info).format = UUID4]",
			Type:       "string",
			problems:   testutils.Problems{{Message: "must specify one of"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_info.proto";

				message Person {
					{{.Type}} {{.FieldName}} = 2 {{.Annotation}};
				}

				message Foo {}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			if diff := test.problems.SetDescriptor(field).Diff(ipAddressFormat.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
