package collection

import "strings"

// SliceContains contains true if the input slice contains the provided substring
func SliceContains(ss []string, subString string) bool {
	for _, s := range ss {
		if strings.Contains(s, subString) {
			return true
		}
	}
	return false
}
