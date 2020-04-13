package util

import (
	"bufio"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

func Uint16(v int) uint16 {
	return uint16(Clamp(v, 0, math.MaxUint16))
}

func Int(v int) int {
	return int(Clamp(v, math.MinInt32, math.MaxInt32))
}

func Int8(v int) int8 {
	return int8(Clamp(v, math.MinInt8, math.MaxInt8))
}

func Uint8(v int) uint8 {
	return uint8(Clamp(v, 0, math.MaxUint8))
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func HasFlagUint8(a, b uint8) bool {
	return a&b != 0
}

func HasFlagUint16(a, b uint16) bool {
	return a&b != 0
}

func StringEqualFold(a string, b ...string) bool {
	for _, v := range b {
		if strings.EqualFold(a, v) {
			return true
		}
	}
	return false
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func MinInt(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// 随机 [low, high]
func RandomInt(low int, high int) int {
	if low == high {
		return low
	}
	if low > high || high == 0 {
		return 0
	}

	return rand.Intn(high-low+1) + low
}

// c# random.next [low, high)
func RandomNext2(low, high int) int {
	return RandomInt(low, high-1)
}

// c# random.next [0, high)
func RandomNext(high int) int {
	return RandomNext2(0, high)
}

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func GetFiles(dir string, allow []string) []string {

	allowMap := map[string]bool{}
	if allow != nil {
		for _, v := range allow {
			allowMap[v] = true
		}
	}

	ret := []string{}
	filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		ext := path.Ext(fpath)
		if allowMap[ext] {
			ret = append(ret, filepath.ToSlash(fpath))
		}

		return nil
	})

	return ret
}

// 按空格拆分字符串。如果加了引号，那么认为是一个字符串
func SplitString(s string) []string {

	ret := []string{}

	start := 0
	var stat rune

	r := []rune(s)

	for i := 0; i < len(r); i++ {
		if unicode.IsSpace(r[i]) {
			if stat == 1 {
				ret = append(ret, string(r[start:i]))
				stat = 0
			}

		} else if r[i] == '\'' || r[i] == '"' {
			if stat == r[i] {
				ret = append(ret, string(r[start:i]))
				stat = 0
			} else {
				if stat == 0 {
					stat = r[i]
					start = i + 1
				}
			}
		} else {
			if stat == 0 {
				start = i
				stat = 1
			}
		}

	}

	if stat != 0 {
		ret = append(ret, string(r[start:]))
	}

	return ret
}

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
		lines = append(lines, RemoveBOM(scanner.Text()))
	}

	return lines
}

// RemoveBOM 删除 windows 保存文件时加入的 bom
func RemoveBOM(s string) string {
	bytes := []byte(s)
	if len(bytes) >= 3 {
		if bytes[0] == 0xef && bytes[1] == 0xbb && bytes[2] == 0xbf {
			return string(bytes[3:])
		}
	}
	return s
}

// FixSeparator 修改路径分隔符
// 比如 SystemScripts\00Default\GuardsHelp.txt
// Windows 下不变
// Unix 下改成 SystemScripts/00Default/GuardsHelp.txt
func FixSeparator(s string) string {
	if string(os.PathSeparator) == "\\" {
		return s
	}
	return strings.Replace(s, "\\", "/", -1)
}
