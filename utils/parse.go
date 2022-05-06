package utils

import (
	"regexp"
	"strconv"
)

func ParseConsensusOutput(target, reg string, matchGr int) float64 {
	match := regexp.MustCompile(reg).FindStringSubmatch(target)
	var result float64
	if match != nil {
		if i, err := strconv.ParseFloat(match[matchGr], 64); err == nil {
			result = i
		}
	}

	return result
}
