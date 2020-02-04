package script

import "reflect"

// 对于if 和 act里面的命令都可以看作一个是一个golang函数
type ScriptFunc struct {
	Name       string
	Func       reflect.Value
	ArgsParser []ArgParseFunc
}

type Function struct {
	Args []reflect.Value
	Func reflect.Value
}

// 用于if，函数要求必须返回bool值
func (c *Function) Check() bool {
	retvars := c.Func.Call(c.Args)
	return retvars[0].Bool()
}

// 用于非if的地方
func (c *Function) Exec() {
	c.Func.Call(c.Args)
}

type Context struct {
	Checks  map[string]*ScriptFunc
	Actions map[string]*ScriptFunc
}

var DefaultContext *Context

var (
	opType reflect.Type
)

func init() {
	opType = reflect.TypeOf(CompareOp(0))

	DefaultContext = &Context{
		Checks:  map[string]*ScriptFunc{},
		Actions: map[string]*ScriptFunc{},
	}
}

func (c *Context) Check(k string, fun interface{}) {
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

	ck := ScriptFunc{
		Name:       k,
		Func:       reflect.ValueOf(fun),
		ArgsParser: checkArgs(typ),
	}

	c.Checks[k] = &ck
}

func checkArgs(funcType reflect.Type) []ArgParseFunc {
	if funcType.NumIn() > 0 {
		parsers := make([]ArgParseFunc, funcType.NumIn())
		for i := 0; i < funcType.NumIn(); i++ {
			argType := funcType.In(i)
			if argType == opType {
				parsers[i] = parseCompare
			} else if argType.Kind() == reflect.String {
				parsers[i] = parseString
			} else if argType.Kind() == reflect.Int {
				parsers[i] = parseInt
			} else {
				panic("not support " + argType.String())
			}
		}

		return parsers
	}

	return []ArgParseFunc{}
}

func Check(k string, fun interface{}) {
	DefaultContext.Check(k, fun)
}

func (c *Context) Action(k string, fun interface{}) {
	typ := reflect.TypeOf(fun)
	if typ.Kind() != reflect.Func {
		panic("func please.")
	}

	ck := ScriptFunc{
		Name:       k,
		Func:       reflect.ValueOf(fun),
		ArgsParser: checkArgs(typ),
	}

	c.Actions[k] = &ck
}

func Action(k string, fun interface{}) {
	DefaultContext.Action(k, fun)
}
