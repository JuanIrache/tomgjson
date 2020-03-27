package tomgjson

import (
	"log"
	"math"
	"strconv"
	"strings"
)

// Like math.Max but with ints
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns both sides of a float number as strings
func sides(n float64) (string, string) {
	sides := strings.Split(strconv.FormatFloat(math.Abs(n), 'f', -1, 64), ".")
	if len(sides) == 1 {
		sides = append(sides, "0")
	}
	if len(sides) != 2 {
		log.Panicf("Badly formatted float: %v %v", n, sides)
	}
	return sides[0], sides[1]
}

// Make sure float values are within mgJSON's valid values

const largestMgjsonNum = 2147483648.0

func validValue(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return math.Max(math.Min(v, largestMgjsonNum), -largestMgjsonNum)
}
