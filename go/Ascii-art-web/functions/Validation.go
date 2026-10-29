package functions

func Validation(word string) bool {
	runes := []rune(word)
	for i := 0; i < len(runes); i++ {
		if !(runes[i] >= 32 && runes[i] <= 126) {
			return false
		}
	}
	return true
}
