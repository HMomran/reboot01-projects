package function

import (
	"strings"
)

func AsciiArt(word, banner string) string {
	var result string

	if !(Validation(word)) {
		return ""
	}

	word = strings.ReplaceAll(word, `\n` , "\n")
	lines := strings.Split(word, "\n")

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
				hold := PrintAscii((((int(lines[j][m])-31)*9 - 7) + k), banner)
				if hold != "" {
					result +=hold
				}else if hold == "" && len(result) < 1 {
					return ""
				}	
			}
			result += "\n"
		}
	}

	return result

}
