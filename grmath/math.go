package grmath

import (
	"fmt"
	"strconv"
)

func Pow(x int, n int) int {
	ans := int(1)
	if n > 0 {
		for n > 0 {
			ans *= x
			n--
		}

	}

	return ans
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
