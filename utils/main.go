package utils

import (
	"sort"
	"strconv"
)

func Sort(mapValue map[string][]string, index int) map[string][]string {
	keys := make([]string, 0, len(mapValue))
	for k := range mapValue {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		vi, _ := strconv.Atoi(mapValue[keys[i]][1])
		vj, _ := strconv.Atoi(mapValue[keys[j]][1])
		return vi > vj
	})

	sortedVSetsResult := make(map[string][]string)
	for _, k := range keys {
		sortedVSetsResult[k] = mapValue[k]
	}

	return sortedVSetsResult
}
