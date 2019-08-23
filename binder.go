package frontpage

import (
	"encoding/json"
	"reflect"

	"github.com/StevenZack/tools/refToolkit"

	"github.com/StevenZack/fasthttp"
)

type binder struct {
	fns  map[string]*Func
	vars *Vars
}

func newBinder(vars *Vars) *binder {
	return &binder{
		fns:  make(map[string]*Func),
		vars: vars,
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
	b.vars.Funcs = append(b.vars.Funcs, *fn)
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
	if len(out) == 0 {
		cx.WriteString("")
		return
	}

	e = checkError(out[len(out)-1])
	if e != nil {
		cx.Error(e.Error(), fasthttp.StatusBadRequest)
		return
	}

	rp, e := json.Marshal(out[0].Interface())
	if e != nil {
		cx.Error(e.Error(), fasthttp.StatusBadRequest)
		return
	}

	cx.SetJsonHeader()
	cx.Write(rp)
}

func checkError(v reflect.Value) error {
	t := v.Type()
	if t.Name() == "error" {
		e := v.Interface()
		if e == nil {
			return nil
		}
		return e.(error)
	}
	return nil
}
