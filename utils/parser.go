package utils

import (
	"strconv"
)

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
