// Copyright 2026 Google LLC
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

package aip0190

import "testing"

func TestIsValidCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"ValidNormal", "CustomerId", true},
		{"ValidXmlHttpRequest", "XmlHttpRequest", true},
		{"ValidIpv6OnIosOptions", "Ipv6OnIosOptions", true},
		{"ValidABTestingService", "ABTestingService", true},
		{"ValidHttpABTestingService", "HttpABTestingService", true},
		// Multiple allowlist entries with trailing lowercase letters.
		{"ValidABTestETaggingRequest", "ABTestETaggingRequest", true},
		// Single letter words can also appear at the end of an allowlisted term.
		{"ValidPlanB", "ExecutePlanB", true},
		// We sometimes require leading and trailing context for single letters .
		{"ValidTypeIIError", "TypeIIErrorResponse", true},
		{"ValidXRay", "XRay", true},
		{"ValidOAuth", "OAuth", true},
		{"ValidABTest", "ABTest", true},
		{"InvalidTShirtService", "TShirtService", false},
		{"InvalidSnakeCase", "snake_case", false},
		{"InvalidLowerCamelCase", "lowerCamelCase", false},
		{"InvalidLowercase", "lowercase", false},
		{"InvalidCustomerID", "CustomerID", false},
		{"InvalidXMLHTTPRequest", "XMLHTTPRequest", false},
		{"InvalidIPv6OnIOSOptions", "IPv6OnIOSOptions", false},
		{"InvalidHTTPABTestingService", "HTTPABTestingService", false},
		{"InvalidEMail", "EMail", false},
		{"InvalidEBook", "EBook", false},
		{"InvalidECommerce", "ECommerce", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isValidCamelCase(test.input)
			if got != test.want {
				t.Errorf("isValidCamelCase(%q) = %t; want %t", test.input, got, test.want)
			}
		})
	}
}
