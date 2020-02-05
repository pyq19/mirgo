package script

import (
	"bufio"
	"container/list"
	"io"
	"os"
	"strings"
	"unicode"
)

func ReadLines(filepath string) (lines []string, err error) {

	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	return ReadLinesByReader(file), nil
}

func ReadLinesByReader(r io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func StartsWithI(str, s string) bool {
	if len(str) < len(s) {
		return false
	}

	return strings.ToUpper(str[:len(s)]) == strings.ToUpper(s)
}

func TrimEnd(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func ListToArray(lst *list.List) []string {
	ret := []string{}

	for it := lst.Front(); it != nil; it = it.Next() {
		ret = append(ret, it.Value.(string))
	}

	return ret
}
