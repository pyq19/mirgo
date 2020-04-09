package script

import (
	"errors"
	"regexp"
	"strings"

	"github.com/yenkeia/mirgo/game/util"
)

var (
	regexInclude = regexp.MustCompile(`#INCLUDE\s*\[([^\n]+)\]\s*(@[^\n]+)`)
	regexInsert  = regexp.MustCompile(`#INSERT\s*\[([^\n]+)\]\s*(@[^\n]+)`)
	regexPage    = regexp.MustCompile(`^(\[[^\n]+\])\s*$`)
)

type PageSource struct {
	Name  string
	Lines []string
}

type PreCompiledScript struct {
	Pages map[string]*PageSource
}

func (p *PreCompiledScript) Add(ps *PageSource) {
	p.Pages[strings.ToUpper(ps.Name)] = ps
}

func (p *PreCompiledScript) Take(name string) *PageSource {
	name = strings.ToUpper(name)
	ps := p.Pages[name]
	delete(p.Pages, name)
	return ps
}

// 按[pagename]拆分脚本行
func precompile(filelines []string) (*PreCompiledScript, error) {
	lines, err := expandScript(filelines)
	if err != nil {
		return nil, err
	}

	var curPage *PageSource

	ret := &PreCompiledScript{}
	ret.Pages = map[string]*PageSource{}

	for _, line := range lines {
		if line[0] == '[' {
			match := regexPage.FindStringSubmatch(line)
			if len(match) > 0 {
				if curPage != nil {
					ret.Add(curPage)
				}

				curPage = &PageSource{
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
		ret.Add(curPage)
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
				match := regexInsert.FindStringSubmatch(line)
				insertLines, err := util.ReadLines(fullpath(match[1]))
				if err != nil {
					return nil, err
				}
				insertLines, err = expandScript(insertLines)
				if err != nil {
					return nil, err
				}
				compiled = append(compiled, insertLines...)
				continue

			} else if StartsWithI(line, "#INCLUDE") {
				match := regexInclude.FindStringSubmatch(line)
				insertLines, err := loadScriptPage(util.FixSeparator(match[1]), match[2])
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
	lines, err := util.ReadLines(fullpath(file))
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
