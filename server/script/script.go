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
)

var EnvirPath = ""

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
	r, err := os.Open(filepath.Join(EnvirPath, file))
	if err != nil {
		return nil, err
	}
	return Load(r)
}

func Load(r io.Reader) (*Script, error) {
	lines := ReadLinesByReader(r)

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

	// buttons := &list.List{}
	// elseButtons := &list.List{}
	// gotoButtons := &list.List{}

	var currentSay = say
	// var currentButtons = buttons

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
				return errors.New("error")
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
		ck, err := ps.parseAction(mp, it.Value.(string))
		if err != nil {
			return nil, err
		}

		actions = append(actions, ck)
	}
	return actions, nil
}

func (ps *PageScript) parseAction(mp map[string]*ScriptFunc, s string) (*Function, error) {
	parts := strings.Split(s, " ")

	method, has := mp[parts[0]]
	if !has {
		return nil, errors.New("no function " + parts[0])
	}

	n := len(method.ArgsParser)
	if n != len(parts)-1 {
		return nil, fmt.Errorf("%s args expect %d got %d", parts[0], n, len(parts)-1)
	}

	inst := Function{}
	inst.Args = make([]reflect.Value, n)
	inst.Func = method.Func

	for i := 0; i < n; i++ {
		value, err := method.ArgsParser[i](parts[i+1])
		if err != nil {
			return nil, err
		}
		inst.Args[i] = value
	}

	return &inst, nil
}

func (sc *Script) Call(page string) ([]string, error) {
	page = strings.ToUpper(page)
	ps, has := sc.Pages[page]
	if !has {
		return nil, errors.New("no page" + page)
	}

	if ps.Check() {
		ps.doActions(ps.ActList)
		return ps.Say, nil
	} else {
		ps.doActions(ps.ElseActList)
		return ps.ElseSay, nil
	}
}

func (ps *PageScript) doActions(actions []*Function) {
	for _, act := range actions {
		act.Exec()
	}
}

func (ps *PageScript) Check() bool {

	if len(ps.CheckList) == 0 {
		return true
	}

	for _, ck := range ps.CheckList {
		if !ck.Check() {
			return false
		}
	}
	return true
}
