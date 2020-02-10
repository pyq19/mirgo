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

// 按空格拆分字符串。如果加了引号，那么认为是一个字符串
func splitString(s string) []string {

	ret := []string{}

	start := 0
	var stat byte

	for i := 0; i < len(s); i++ {
		// fmt.Println("[" + fmt.Sprintf("%c-%c", stat, s[i]) + "]")
		switch s[i] {
		case ' ', '\'', '"':
			if stat == s[i] {
				ret = append(ret, s[start:i])
				stat = 0
			} else {
				if stat == 0 {
					stat = s[i]
					start = i + 1
				}
			}
		default:
			if stat == 0 {
				stat = ' '
				start = i
			}
		}
	}

	if stat != 0 {
		ret = append(ret, s[start:])
	}

	return ret
}
