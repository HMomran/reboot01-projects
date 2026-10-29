package function

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func PrintAscii(num int, banner string) string {
	bannerAlternate := map[string]string{
		"shadow":     "https://raw.githubusercontent.com/zsandibe/ascii-art/master/shadow.txt",
		"standard":   "https://raw.githubusercontent.com/zsandibe/ascii-art/master/standard.txt",
		"thinkertoy": "https://raw.githubusercontent.com/zsandibe/ascii-art/master/thinkertoy.txt",
	}

	var word string
	banner = strings.ToLower(banner)
	var reader io.ReadCloser
	var err error

	if banner == "shadow" {
		reader, err = os.Open("banner/shadow.txt")
	} else if banner == "standard" {
		reader, err = os.Open("banner/standard.txt")
	} else if banner == "thinkertoy" {
		reader, err = os.Open("banner/ThinkerToy.txt")
	} else {
		fmt.Println("Unknown banner style")
		return ""
	}

	if err != nil {
		resp, err := http.Get(bannerAlternate[banner])
		if err != nil {
			fmt.Println("Failed to fetch from GitHub:", err)
			return ""
		}
		reader = resp.Body 
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	start := 1
	for scanner.Scan() {
		if start == num {
			word += scanner.Text()
			break
		}
		start++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	return word
}
