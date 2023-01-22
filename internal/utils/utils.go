package utils

import "math"

func IntMax(a int, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
