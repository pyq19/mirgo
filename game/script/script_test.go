package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/yenkeia/mirgo/game/util"
)

func Print(arr []string) {

	for i, v := range arr {
		fmt.Println(i, "[["+v+"]]")
	}
}

func PrintJson(v interface{}) {
	j, _ := json.MarshalIndent(v, "", "\t")
	fmt.Print(string(j))
}

func setpath() {
	gopath := os.Getenv("GOPATH")
	SearchPaths = []string{filepath.Join(gopath, "src/github.com/yenkeia/mirgo/dotnettools/database/Envir")}
}

func TestXXX(X *testing.T) {
	fmt.Println(StartsWithI("#INSERT [", "#INSERT"))
	fmt.Println(StartsWithI("#INSerT [", "#INSERT"))
	fmt.Println(StartsWithI("#INSer", "#INSERT"))
	fmt.Println(StartsWithI("[@MAIN]", "[@Main]"))

	v := regexPage.FindStringSubmatch("[@main]")
	Print(v)

	a := regexInclude.FindStringSubmatch("#INCLUDE [SystemScripts/SharedNPCS/Tavern.txt] @Main")
	Print(a)
}

func TestPrecompile(t *testing.T) {
	setpath()

	v, err := loadScriptPage("SystemScripts/SharedNPCS/Tavern.txt", "@Main")
	if err != nil {
		panic(err)
	}
	Print(v)
	lines, err := util.ReadLines(fullpath("NPCs/BichonProvince/BichonWall/Sir.MoguBW.txt"))
	if err != nil {
		panic(err)
	}

	m, err := precompile(lines)

	if err != nil {
		panic(err)
	}
	fmt.Print(m)
	PrintJson(m)
}

func TestLogic(t *testing.T) {
	var test = `[@MAIN]
#IF
CHECKPKPOINT > 2
#SAY
I will not help an evil person like you...
	
	
<Close/@exit>
#ACT
TEST 
TEST 1
TEST 1 2
TEST2 100
#ELSEACT
GOTO @Main-1

[@Main-1]
#SAY
Welcome, what can I do for you?
	
<Sell/@Sell> Meat.
<Ask/@Meathelp> about how to gain meat.
	
<Close/@Exit>
`

	// skip uint32
	DefaultContext.AddParser(reflect.TypeOf(uint32(0)), nil)

	Check("CHECKPKPOINT", func(op CompareOp, v int) bool {
		return CompareInt(op, 3, v)
	})

	Action("TEST", func(v1, v2 int) {
		fmt.Println(v1, v2)
	}, -909, -8080)

	Action("TEST2", func(v1 uint32, v2 int) {
		fmt.Println(v1, v2)
	})

	sc, err := Load(bytes.NewReader([]byte(test)))
	if err != nil {
		panic(err)
	}

	say, err := sc.Call("[@main]", uint32(1009))
	if err != nil {
		fmt.Println("err", err)
	}
	Print(say)
}

func replaces(s string) string {
	fmt.Println(s)
	return ""
}
