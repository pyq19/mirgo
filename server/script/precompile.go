package script

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regexInclude = regexp.MustCompile(`#INCLUDE\s*\[([^\n]+)\]\s*(@[^\n]+)`)
	regexPage    = regexp.MustCompile(`^\[([^\n]+)\]\s*$`)
)

type PageScript struct {
	Name  string
	Lines []string
}

type PreCompiledScript struct {
	Pages []*PageScript
}

func (p *PreCompiledScript) GetPage(name string) *PageScript {
	for _, v := range p.Pages {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return v
		}
	}
	return nil
}

// 按[pagename]拆分脚本行
func precompile(filelines []string) (*PreCompiledScript, error) {
	lines, err := expandScript(filelines)
	if err != nil {
		return nil, err
	}

	var curPage *PageScript

	ret := &PreCompiledScript{}
	ret.Pages = []*PageScript{}

	for _, line := range lines {
		if line[0] == '[' {
			match := regexPage.FindStringSubmatch(line)
			if len(match) > 0 {
				if curPage != nil {
					ret.Pages = append(ret.Pages, curPage)
				}

				curPage = &PageScript{
					Name:  match[1],
					Lines: []string{},
				}
				continue
			}
		}

		if curPage != nil {
			curPage.Lines = append(curPage.Lines, line)
		} else {
			// skip
			// fmt.Println("skip=>", line)
		}
	}

	if curPage != nil {
		ret.Pages = append(ret.Pages, curPage)
	}

	return ret, nil
}

func skipLine(line string) bool {
	//
	if line == "" {
		return true
	}
	// 注释
	if line[0] == ';' {
		return true
	}
	return false
}

// 去掉注释和空行，并处理include
func expandScript(lines []string) ([]string, error) {
	compiled := []string{}

	for _, line := range lines {
		if skipLine(line) {
			continue
		}
		if line[0] == '#' {
			if StartsWithI(line, "#INSERT") {
				panic("#INSERT not impl yet.")

			} else if StartsWithI(line, "#INCLUDE") {
				match := regexInclude.FindStringSubmatch(line)
				insertLines, err := loadScriptPage(match[1], match[2])
				if err != nil {
					return nil, err
				}
				compiled = append(compiled, insertLines...)
				continue
			}
		}

		compiled = append(compiled, line)
	}

	return compiled, nil
}

func loadScriptPage(file, page string) ([]string, error) {
	lines, err := ReadLines(filepath.Join(EnvirPath, file))
	if err != nil {
		return nil, err
	}

	page = "[" + page + "]"

	stat := 0

	ret := []string{}

	for _, line := range lines {
		if skipLine(line) {
			continue
		}
		switch stat {
		case 0:
			if line[0] == '[' && StartsWithI(line, page) {
				stat = 1
			}

		case 1:
			if line[0] == '{' {
				stat = 2
			}
		case 2:
			if line[0] == '}' {
				return ret, nil
			}

			ret = append(ret, line)
		}
	}

	return nil, errors.New("sytax error:" + file)
}
