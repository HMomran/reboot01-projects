package function

func Validation(word string) bool {
	test := true
	array := []rune(word)
	for i := 0; i < len(array); i++ {
		if !(array[i] >= 32 && array[i] <= 126) {
			test = false
		}
	}

	return test
}
