package lint

// aliasMap stores legacy names for some rules.
var aliasMap = map[string]string{
	"core::0140::lower-snake":                  "naming-format",
	"core::0140::enum-names::abbreviations":    "abbreviations",
	"core::0140::field-names::abbreviations":   "abbreviations",
	"core::0140::message-names::abbreviations": "abbreviations",
	"core::0140::method-names::abbreviations":  "abbreviations",
	"core::0140::service-names::abbreviations": "abbreviations",
}
