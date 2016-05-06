package beater

import "strings"

//GREEN cluster status
const GREEN = 2

//YELLOW cluster status
const YELLOW = 1

//RED cluster status
const RED = 0

//UNKNOWN cluster status
const UNKNOWN = -1

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func GetNumericalClusterStatus(s string) int64 {

	// numeric interpretations of cluster health status
	if strings.EqualFold(s, "green") {
		return GREEN
	} else if strings.EqualFold(s, "yellow") {
		return YELLOW
	} else if strings.EqualFold(s, "red") {
		return RED
	}
	return UNKNOWN
}
