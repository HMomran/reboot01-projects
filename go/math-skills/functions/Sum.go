package functions

func Sum(array []float64) float64 {
	var sum float64
	for i := 0; i < len(array); i++ {
		sum += array[i]
	}
	return sum
}
