// Package utils provide helped util functions
package utils

func ContainsAll[C comparable](a, b []C) bool {
	set := make(map[C]struct{})

	for _, v := range b {
		set[v] = struct{}{}
	}

	for _, v := range a {
		if _, found := set[v]; !found {
			return false
		}
	}

	return true
}
