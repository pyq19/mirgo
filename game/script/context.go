package script

import (
	"reflect"
)

// 对于if 和 act里面的命令都可以看作一个是一个golang函数
type ScriptFunc struct {
	Name       string
	Func       reflect.Value
	ArgsParser []*ArgParser
	OptionArgs []reflect.Value
}

type Function struct {
	Skiped bool
	Args   []reflect.Value
	Func   reflect.Value
}

type CMDBreak struct {
}

type CMDGoto struct {
	GOTO string
}

func (c *Function) doExec(args ...interface{}) []reflect.Value {
	if c.Skiped {
		for i, v := range args {
			c.Args[i] = reflect.ValueOf(v)
		}
	}
	return c.Func.Call(c.Args)
}

// 用于if，函数要求必须返回bool值
func (c *Function) Check(args ...interface{}) bool {
	return c.doExec(args...)[0].Bool()
}

// 用于非if的地方
func (c *Function) Exec(args ...interface{}) interface{} {
	ret := c.doExec(args...)

	if ret != nil && len(ret) > 0 {
		return ret[0].Interface()
	}

	return nil
}

type Context struct {
	Checks  map[string]*ScriptFunc
	Actions map[string]*ScriptFunc
	parsers map[reflect.Type]*ArgParser
}

var DefaultContext *Context

func NewContext() *Context {
	r := &Context{
		Checks:  map[string]*ScriptFunc{},
		Actions: map[string]*ScriptFunc{},
		parsers: map[reflect.Type]*ArgParser{},
	}

	r.AddParser(reflect.TypeOf(bool(false)), ParseBool)
	r.AddParser(reflect.TypeOf(int(0)), ParseInt)
	r.AddParser(reflect.TypeOf(string("")), ParseString)
	r.AddParser(reflect.TypeOf(GT), ParseCompare)

	return r
}

// 直接运行某字符串。如: Move map001 100 100
func (c *Context) Exec(s string, args ...interface{}) (interface{}, error) {
	f, err := parseAction(c.Actions, s)
	if err != nil {
		return nil, err
	}

	return f.Exec(args...), nil
}

// f==nil 表示不解析该参数，从外部传入。
func (c *Context) AddParser(typ reflect.Type, f ArgParseFunc) {
	c.parsers[typ] = &ArgParser{Fun: f}
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
	f := &ScriptFunc{
		Name:       k,
		Func:       reflect.ValueOf(fun),
		ArgsParser: c.checkArgs(reflect.TypeOf(fun)),
		OptionArgs: c.checkOptions(options...),
	}

	return f
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

func (c *Context) checkArgs(funcType reflect.Type) []*ArgParser {
	n := funcType.NumIn()
	if n > 0 {
		parsers := make([]*ArgParser, n)
		for i := 0; i < n; i++ {
			argType := funcType.In(i)

			par, has := c.parsers[argType]
			if !has {
				panic("not support arguments type " + argType.String())
			}
			parsers[i] = par
		}

		return parsers
	}

	return []*ArgParser{}
}

func init() {
	DefaultContext = NewContext()

	Action("GOTO", _GOTO)
	Action("BREAK", _BREAK)
}

func _GOTO(page string) CMDGoto {
	return CMDGoto{GOTO: "[" + page + "]"}
}
func _BREAK() CMDBreak {
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
