package functions

import "math"

func Variance(array []float64) float64 {
	mean := Average(array)
	var variance float64
	for i := 0; i < len(array); i++ {
		variance += math.Pow(array[i]-mean, 2)
	}
	variance = variance / float64(len(array))
	return variance

}
