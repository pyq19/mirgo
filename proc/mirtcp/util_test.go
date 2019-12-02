package mirtcp

import (
	"strconv"
	"strings"
	"testing"
)

func _String(bytes []byte) string {
	var strSlice = []string{}
	for _, b := range bytes {
		strSlice = append(strSlice, strconv.Itoa(int(b)))
	}
	res := strings.Join(strSlice, ", ")
	return "[" + res + "]"
}


func TestString(t *testing.T) {
	bytes := []byte{111, 11, 1}
	t.Log(_String(bytes))
}
