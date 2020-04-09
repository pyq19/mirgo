package script

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/yenkeia/mirgo/game/util"
)

var SearchPaths = []string{}

func fullpath(file string) string {
	file = strings.ReplaceAll(file, "\\", "/")

	if filepath.IsAbs(file) {
		return file
	}
	for _, v := range SearchPaths {
		full := filepath.Join(v, file)
		if util.IsFile(full) {
			return full
		}
	}
	return file
}

type Script struct {
	Types  []int
	Quests []int
	Goods  []string
	Pages  map[string]*PageScript
}

type PageScript struct {
	Name        string
	CheckList   []*Function
	ActList     []*Function
	ElseActList []*Function
	Say         []string
	ElseSay     []string
}

func LoadFile(file string) (*Script, error) {
	r, err := os.Open(fullpath(file))
	if err != nil {
		return nil, errors.New("LoadFile error " + file + ":" + err.Error())
	}
	return Load(r)
}

func Load(r io.Reader) (*Script, error) {
	lines := util.ReadLinesByReader(r)

	sc := &Script{}

	obj, err := precompile(lines)
	if err != nil {
		return nil, err
	}

	if err := sc.parseGoods(obj.Take("[Trade]")); err != nil {
		return nil, err
	}
	if err := sc.parseTypes(obj.Take("[Types]")); err != nil {
		return nil, err
	}
	if err := sc.parseQuests(obj.Take("[Quests]")); err != nil {
		return nil, err
	}

	sc.Pages = map[string]*PageScript{}
	for _, ps := range obj.Pages {
		page := &PageScript{Name: ps.Name}

		if err := page.parsePage(ps); err != nil {
			return nil, err
		}

		sc.Pages[strings.ToUpper(page.Name)] = page
	}

	return sc, nil
}

func (sc *Script) parseIntArray(p *PageSource) ([]int, error) {
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

func (sc *Script) parseTypes(p *PageSource) (err error) {
	sc.Types, err = sc.parseIntArray(p)
	return
}

func (sc *Script) parseQuests(p *PageSource) (err error) {
	sc.Quests, err = sc.parseIntArray(p)
	return
}

func (sc *Script) parseGoods(p *PageSource) error {
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

var (
	regexSharp = regexp.MustCompile(`#(\w+)`)
)

func (ps *PageScript) parsePage(p *PageSource) error {
	checks := &list.List{}
	acts := &list.List{}
	say := &list.List{}
	elseActs := &list.List{}
	elseSay := &list.List{}

	var currentSay = say

	for i := 0; i < len(p.Lines); i++ {
		line := p.Lines[i]

		if line[0] == '#' {
			match := regexSharp.FindStringSubmatch(line)

			switch strings.ToUpper(match[1]) {
			case "IF":
				currentSay = checks
			case "SAY":
				currentSay = say
			case "ACT":
				currentSay = acts
			case "ELSEACT":
				currentSay = elseActs
			case "ELSESAY":
				currentSay = elseSay
			default:
				return errors.New("error:" + p.Name + "---" + match[1])
			}
			continue
		}

		currentSay.PushBack(TrimEnd(line))
	}
	ps.Say = ListToArray(say)
	ps.ElseSay = ListToArray(elseSay)

	var err error

	ps.CheckList, err = ps.parseActions(DefaultContext.Checks, checks)
	if err != nil {
		return err
	}

	ps.ActList, err = ps.parseActions(DefaultContext.Actions, acts)
	if err != nil {
		return err
	}

	ps.ElseActList, err = ps.parseActions(DefaultContext.Actions, elseActs)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PageScript) parseActions(mp map[string]*ScriptFunc, lst *list.List) ([]*Function, error) {
	actions := []*Function{}

	for it := lst.Front(); it != nil; it = it.Next() {
		ck, err := parseAction(mp, it.Value.(string))
		if err != nil {
			return nil, err
		}

		actions = append(actions, ck)
	}
	return actions, nil
}

func parseAction(mp map[string]*ScriptFunc, s string) (*Function, error) {
	parts := util.SplitString(s)

	funName := strings.ToUpper(parts[0])

	method, has := mp[funName]
	if !has {
		return nil, errors.New("no function [" + funName + "]")
	}

	argsSkip := 0
	for _, parser := range method.ArgsParser {
		if parser.Fun == nil {
			argsSkip++
		} else {
			break
		}
	}
	expect := len(method.ArgsParser) - argsSkip

	got := len(parts) - 1
	opt := 0
	if method.OptionArgs != nil {
		opt = len(method.OptionArgs)
	}

	if expect != got && expect > got+opt {
		return nil, fmt.Errorf("%s args expect %d got %d", funName, expect, got)
	}

	inst := Function{}
	inst.Args = make([]reflect.Value, expect+argsSkip)
	inst.Func = method.Func

	if argsSkip > 0 {
		inst.Skiped = true
	}

	for i := 0; i < expect; i++ {
		if i >= got {
			inst.Args[argsSkip+i] = method.OptionArgs[i-(expect-opt)]
		} else {
			value, err := method.ArgsParser[argsSkip+i].Fun(parts[i+1])
			if err != nil {
				return nil, err
			}
			inst.Args[argsSkip+i] = value
		}
	}

	return &inst, nil
}

// call
func (sc *Script) Call(page string, args ...interface{}) ([]string, error) {

	page = strings.ToUpper(page)
	ps, has := sc.Pages[page]
	if !has {
		return nil, errors.New("no page" + page)
	}

	var acts []*Function
	var say []string

	if ps.Check(args...) {
		acts = ps.ActList
		say = ps.Say
	} else {
		acts = ps.ElseActList
		say = ps.ElseSay
	}

	for _, act := range acts {
		cmd := act.Exec(args...)
		shouldBreak := false
		if cmd != nil {
			switch cmd.(type) {
			case CMDBreak:
				shouldBreak = true
			case CMDGoto:
				return sc.Call(cmd.(CMDGoto).GOTO, args...)
			}
		}
		if shouldBreak {
			break
		}
	}

	return say, nil
}

func (ps *PageScript) Check(args ...interface{}) bool {

	if len(ps.CheckList) == 0 {
		return true
	}

	for _, ck := range ps.CheckList {
		if !ck.Check(args...) {
			return false
		}
	}
	return true
}
