package frontpage

import (
	"reflect"

	"github.com/StevenZack/frontpage/util"

	"github.com/StevenZack/tools/refToolkit"

	"github.com/StevenZack/fasthttp"
)

type binder struct {
	fns map[string]*Func
}

func newBinder() *binder {
	return &binder{
		fns: make(map[string]*Func),
	}
}

func (b *binder) bind(v interface{}) {
	name, e := refToolkit.GetFuncName(v)
	if e != nil {
		panic(e)
	}
	_, ok := b.fns[name]
	if ok {
		panic("Bind():funcname '" + name + "' duplicated")
	}

	fn, e := NewFunc(v)
	if e != nil {
		panic(e)
	}

	b.fns[name] = fn
}
func (b *binder) handleCall(cx *fasthttp.RequestCtx) {
	funcName := cx.GetPathParam("funcName")
	if funcName == "" {
		cx.Error("funcName not found", fasthttp.StatusBadRequest)
		return
	}

	fn, ok := b.fns[funcName]
	if !ok {
		cx.Error("funcName '"+funcName+"' not bound", fasthttp.StatusBadRequest)
		return
	}

	in, e := fn.ParseIn(string(cx.PostBody()))
	if e != nil {
		cx.Error(e.Error(), fasthttp.StatusBadRequest)
		return
	}

	out := reflect.ValueOf(fn.i).Call(in)
	str, e := util.RefVsToJson(out)
	if e != nil {
		cx.Error(e.Error(), fasthttp.StatusBadRequest)
		return
	}

	cx.SetJsonHeader()
	cx.WriteString(str)
}
