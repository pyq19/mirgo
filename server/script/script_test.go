package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
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
	EnvirPath = filepath.Join(gopath, "src/github.com/yenkeia/mirgo/dotnettools/database/Envir")
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
	lines, err := ReadLines(filepath.Join(EnvirPath, "NPCs/BichonProvince/BichonWall/Sir.MoguBW.txt"))
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
#ELSEACT
GOTO @Main-1

[@Main-1]
#SAY
Welcome, what can I do for you?
	
<Sell/@Sell> Meat.
<Ask/@Meathelp> about how to gain meat.
	
<Close/@Exit>
`

	Check("CHECKPKPOINT", func(npc, player interface{}, op CompareOp, v int) bool {
		return CompareInt(op, 3, v)
	})

	Action("TEST", func(npc, player interface{}, v1, v2 int) {
		fmt.Println(v1, v2)
	}, -909, -8080)

	sc, err := Load(bytes.NewReader([]byte(test)))
	if err != nil {
		panic(err)
	}

	say, err := sc.Call(1, 1, "[@main]")
	if err != nil {
		fmt.Println("err", err)
	}
	Print(say)
}

func TestSplit(t *testing.T) {
	Print(splitString("hello world += 1 "))
	Print(splitString("\"hello world += 1 \""))
	Print(splitString(`"hello world" += 1 "123"`))
	Print(splitString(`CHECKGOLD > 100`))
	Print(splitString(`CHECKQUEST 154 COMPLETE`))
	s := "[234]"
	fmt.Println(s[1 : len(s)-1])

	var regNPCHotkey = regexp.MustCompile(`\<\$\w+\>`)

	fmt.Println(regNPCHotkey.ReplaceAllStringFunc("asdjd<$username><hello>", replaces))
}

func replaces(s string) string {
	fmt.Println(s)
	return ""
}
