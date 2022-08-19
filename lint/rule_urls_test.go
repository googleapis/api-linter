package lint

import (
	"testing"
)

func TestCoreRuleURL(t *testing.T) {
	tests := []struct {
		name string
		rule string
		url  string
	}{
		{"CoreRule", "core::0122::camel-case-uris", "https://linter.aip.dev/122/camel-case-uris"},
		{"NotCoreRule", "test::0122::camel-case-uris", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := coreRuleURL(test.rule); got != test.url {
				t.Errorf("coreRuleURL(%s) got %s, but want %s", test.name, got, test.url)
			}
		})
	}
}

func TestClientLibrariesRuleURL(t *testing.T) {
	tests := []struct {
		name string
		rule string
		url  string
	}{
		{"ClientLibrariesRule", "client-libraries::4232::repeated-fields", "https://linter.aip.dev/4232/repeated-fields"},
		{"NotClientLibrariesRule", "test::0122::camel-case-uris", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := clientLibrariesRuleUrl(test.rule); got != test.url {
				t.Errorf("clientLibrariesRuleUrl(%s) got %s, but want %s", test.name, got, test.url)
			}
		})
	}
}

func TestCloudRuleURL(t *testing.T) {
	tests := []struct {
		name string
		rule string
		url  string
	}{
		{"CloudRule", "cloud::2500::generic-fields", "https://linter.aip.dev/2500/generic-fields"},
		{"NotCloudRule", "test::0122::camel-case-uris", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := cloudRuleUrl(test.rule); got != test.url {
				t.Errorf("cloudRuleUrl(%s) got %s, but want %s", test.name, got, test.url)
			}
		})
	}
}

func TestGetRuleURL(t *testing.T) {
	mapping1 := func(name string) string {
		if name == "one" {
			return "ONE"
		}
		return ""
	}
	mapping2 := func(name string) string {
		if name == "two" {
			return "TWO"
		}
		return ""
	}
	ruleURLMappings := []func(string) string{mapping1, mapping2}

	tests := []struct {
		name     string
		ruleName string
		ruleURL  string
	}{
		{"MappingOne", "one", "ONE"},
		{"MappingTwo", "two", "TWO"},
		{"NoMapping", "zero", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getRuleURL(test.ruleName, ruleURLMappings); got != test.ruleURL {
				t.Errorf("getRuleURL(%s) got %s, but want %s", test.ruleName, got, test.ruleURL)
			}
		})
	}
}
