package lint

import "fmt"

// A list of functions, each of which returns the group name for the given AIP
// number and if no group is found, returns an empty string.
// NOTE: the list will be evaluated in the FILO order.
//
// At Google, we inject additional group naming functions into this list.
// Example: google_aip_groups.go
// package lint
// func init() {
//   aipGroups = append(aipGroups, aipInternalGroup)
// }
//
// func aipInternalGroup(aip int) string {
//   if aip > 9000 {
//	   return "internal"
//   }
//   return ""
// }
var aipGroups = []func(int) string{
	aipCoreGroup,
}

func aipCoreGroup(aip int) string {
	if aip > 0 && aip < 1000 {
		return "core"
	}
	return ""
}

// getRuleGroup takes an AIP number and returns the appropriate group.
// It panics if no group is found.
func getRuleGroup(aip int, groups []func(int) string) string {
	for i := len(groups) - 1; i >= 0; i-- {
		if group := groups[i](aip); group != "" {
			return group
		}
	}
	panic(fmt.Sprintf("Invalid AIP number %d: no available group.", aip))
}
