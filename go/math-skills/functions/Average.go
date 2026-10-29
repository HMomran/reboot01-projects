package functions

func Average(array []float64) float64 {
	sum := Sum(array)
	return float64(sum) / float64(len(array))
}
