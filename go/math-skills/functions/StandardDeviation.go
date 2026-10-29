package functions

import "math"

func StandardDeviation(array []float64) float64 {
	return math.Sqrt(Variance(array))
}
