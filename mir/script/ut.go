package script

import (
	"container/list"
	"strings"
	"unicode"
)

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
