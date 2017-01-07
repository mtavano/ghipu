package utils

import "sort"

func sortKeys(m map[string]string) []string {
	sortedKeys := make([]string, len(m))

	i := 0
	for k := range m {
		sortedKeys[i] = k
		i++
	}

	sort.Strings(sortedKeys)

	return sortedKeys
}
