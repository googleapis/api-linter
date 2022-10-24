package lint

// aliasMap stores legacy names for some rules.
// At Google, we inject rule-alias mapping into this map.
// Example:
// We will compile an addition file -- "google_rule_aliases.go".
// ````````````````````````````````````````````````````````````
// package lint
//
//	func init() {
//		aliasMap["core::0140::lower-snake"] = "naming-format"
//		aliasMap["core::0140::enum-names::abbreviations"] = "abbreviations"
//	}
//
// ````````````````````````````````````````````````````````````
var aliasMap = map[string]string{}
