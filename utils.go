package tomgjson

import "log"

func check(e error) {
	if e != nil {
		log.Panic("Error:", e)
	}
}

func mMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
