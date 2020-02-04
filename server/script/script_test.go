package script

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	// v, err := loadScriptPage("SystemScripts/SharedNPCS/Tavern.txt", "@Main")
	// if err != nil {
	// 	panic(err)
	// }
	// Print(v)
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
