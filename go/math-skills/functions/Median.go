package functions

import "sort"

func Median(array []float64) float64 {

	sort.Float64s(array)
	if len(array)%2 == 1 {
		return array[len(array)/2]
	} else {
		num1 := array[len(array)/2]
		num2 := array[(len(array)/2)-1]
		return (num1 + num2) / 2.0
	}
}
