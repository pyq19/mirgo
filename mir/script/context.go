package script

import "reflect"

// 对于if 和 act里面的命令都可以看作一个是一个golang函数
type ScriptFunc struct {
	Name       string
	Func       reflect.Value
	ArgsParser []ArgParseFunc
	OptionArgs []reflect.Value
}

type Function struct {
	Args []reflect.Value
	Func reflect.Value
}

type CMDBreak struct {
}

type CMDGoto struct {
	GOTO string
}

// 用于if，函数要求必须返回bool值
func (c *Function) Check(npc, player interface{}) bool {
	c.Args[0] = reflect.ValueOf(npc)
	c.Args[1] = reflect.ValueOf(player)
	retvars := c.Func.Call(c.Args)
	return retvars[0].Bool()
}

// 用于非if的地方
func (c *Function) Exec(npc, player interface{}) interface{} {
	c.Args[0] = reflect.ValueOf(npc)
	c.Args[1] = reflect.ValueOf(player)
	ret := c.Func.Call(c.Args)
	if ret != nil && len(ret) > 0 {
		return ret[0].Interface()
	}
	return nil
}

type Context struct {
	Checks  map[string]*ScriptFunc
	Actions map[string]*ScriptFunc
	parsers map[reflect.Type]ArgParseFunc
}

var DefaultContext *Context

var (
	opType reflect.Type
)

func _GOTO(a, b interface{}, page string) CMDGoto {
	return CMDGoto{GOTO: "[" + page + "]"}
}
func _BREAK(a, b interface{}) CMDBreak {
	return CMDBreak{}
}

// 给类型添加解析函数
func AddParser(typ reflect.Type, f ArgParseFunc) {
	DefaultContext.AddParser(typ, f)
}

// 函数名，函数，可选参数默认值
func Check(k string, fun interface{}, options ...interface{}) {
	DefaultContext.Check(k, fun, options...)
}

// 函数名，函数，可选参数默认值
func Action(k string, fun interface{}, options ...interface{}) {
	DefaultContext.Action(k, fun, options...)
}

func (c *Context) AddParser(typ reflect.Type, f ArgParseFunc) {
	c.parsers[typ] = f
}

func (c *Context) Check(k string, fun interface{}, options ...interface{}) {
	typ := reflect.TypeOf(fun)
	if typ.Kind() != reflect.Func {
		panic("func please.")
	}

	if typ.NumOut() != 1 {
		panic("check func should return bool")
	}

	out0 := typ.Out(0)
	if out0.Kind() != reflect.Bool {
		panic("check func should return bool")
	}

	c.Checks[k] = c.makefunc(k, fun, options...)
}

func (c *Context) Action(k string, fun interface{}, options ...interface{}) {
	typ := reflect.TypeOf(fun)
	if typ.Kind() != reflect.Func {
		panic("func please.")
	}

	c.Actions[k] = c.makefunc(k, fun, options...)
}

func (c *Context) makefunc(k string, fun interface{}, options ...interface{}) *ScriptFunc {
	return &ScriptFunc{
		Name:       k,
		Func:       reflect.ValueOf(fun),
		ArgsParser: c.checkArgs(reflect.TypeOf(fun)),
		OptionArgs: c.checkOptions(options...),
	}
}

func (c *Context) checkOptions(options ...interface{}) []reflect.Value {
	if options != nil && len(options) > 0 {
		optargs := make([]reflect.Value, len(options))
		for i, o := range options {
			optargs[i] = reflect.ValueOf(o)
		}
		return optargs
	}
	return nil
}

func (c *Context) checkArgs(funcType reflect.Type) []ArgParseFunc {
	n := funcType.NumIn() - argsSkip
	if n > 0 {
		parsers := make([]ArgParseFunc, n)
		for i := 0; i < n; i++ {
			argType := funcType.In(i + argsSkip)

			par, has := c.parsers[argType]
			if !has {
				panic("not support arguments type " + argType.String())
			}
			parsers[i] = par
		}

		return parsers
	}

	return []ArgParseFunc{}
}

func init() {
	opType = reflect.TypeOf(CompareOp(0))

	DefaultContext = &Context{
		Checks:  map[string]*ScriptFunc{},
		Actions: map[string]*ScriptFunc{},
		parsers: map[reflect.Type]ArgParseFunc{},
	}

	AddParser(reflect.TypeOf(bool(false)), ParseBool)
	AddParser(reflect.TypeOf(int(0)), ParseInt)
	AddParser(reflect.TypeOf(string("")), ParseString)
	AddParser(reflect.TypeOf(GT), ParseCompare)

	Action("GOTO", _GOTO)
	Action("BREAK", _BREAK)
}
