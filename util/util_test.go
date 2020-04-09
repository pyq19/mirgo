package util

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {

	test := func(s string) {
		arr := SplitString(s)
		fmt.Println("---test[" + s + "]")
		for i, v := range arr {
			fmt.Println(i, "["+v+"]")
		}
	}
	test("A               B c")
	test("1/30   EnergyRepulsor")
	test("1/30   \"EnergyRepu ls  or\"")
}
