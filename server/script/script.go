package script

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var EnvirPath = ""

type Script struct {
	Types  []int
	Quests []int
	Goods  []string
}

func LoadFile(file string) error {
	lines, err := ReadLines(filepath.Join(EnvirPath, file))
	if err != nil {
		return err
	}

	sc := &Script{}

	obj, err := precompile(lines)
	if err != nil {
		return err
	}

	if err := sc.parseMain(obj.GetPage("@Main")); err != nil {
		return err
	}
	if err := sc.parseGoods(obj.GetPage("@Trade")); err != nil {
		return err
	}
	if err := sc.parseTypes(obj.GetPage("@Types")); err != nil {
		return err
	}
	if err := sc.parseQuests(obj.GetPage("@Quests")); err != nil {
		return err
	}

	return nil
}

func (sc *Script) parseIntArray(p *PageScript) ([]int, error) {
	ret := []int{}

	if p == nil {
		return ret, nil
	}

	for _, v := range p.Lines {
		id, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		ret = append(ret, id)
	}

	return ret, nil
}

func (sc *Script) parseTypes(p *PageScript) (err error) {
	sc.Types, err = sc.parseIntArray(p)
	return
}

func (sc *Script) parseQuests(p *PageScript) (err error) {
	sc.Quests, err = sc.parseIntArray(p)
	return
}

func (sc *Script) parseGoods(p *PageScript) error {
	sc.Goods = []string{}

	if p == nil {
		return nil
	}

	for _, v := range p.Lines {
		goodsName := strings.TrimSpace(v)
		sc.Goods = append(sc.Goods, goodsName)
	}

	return nil
}

var regexSharp = regexp.MustCompile(`#(\w+)`)

func (sc *Script) parseMain(p *PageScript) error {
	// TODO

	for _, line := range p.Lines {
		if line[0] == '#' {
			match := regexSharp.FindStringSubmatch(line)

			switch strings.ToUpper(match[1]) {
			case "IF":
			case "ELSEACT":
			case "ELSESAY":
			case "SAY":
			case "ACT":
			default:
				return errors.New("error")
			}
		}
	}

	return nil
}
