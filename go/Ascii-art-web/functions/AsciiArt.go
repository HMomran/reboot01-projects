package functions

import (
	"strings"
)

func AsciiArt(word, banner string) string {
	var result string

	if len(word) > 50 {
		return ""
	}

	if !(Validation(word)) {
		return ""
	}

	processed := strings.ReplaceAll(word, `\n`, "\n")
	lines := strings.Split(processed, "\n")

	for j := 0; j < len(lines); j++ {
		if lines[j] == "" && j == len(lines)-1 {
			break
		}

		if lines[j] == "" {
			result += "\n"
			continue
		}

		for k := 0; k < 8; k++ {
			for m := 0; m < len(lines[j]); m++ {
				hold := AsciiPrint((((int(lines[j][m])-31)*9 - 7) + k), banner)
				if hold != "" {
					result += hold
				} else {
					return ""
				}

			}
			result += "\n"
		}
	}

	return result

}
