package script

import (
	"bufio"
	"os"
	"strings"
)

func ReadLines(filepath string) (lines []string, err error) {

	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}

func StartsWithI(str, s string) bool {
	if len(str) < len(s) {
		return false
	}

	return strings.ToUpper(str[:len(s)]) == strings.ToUpper(s)
}
