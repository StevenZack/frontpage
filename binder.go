package frontpage

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/StevenZack/frontpage/util"
	"github.com/StevenZack/mux"
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
	name, e := util.GetFuncName(v)
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
func (b *binder) handleCall(w http.ResponseWriter, r *http.Request) {
	funcName, e := mux.GetURIParam(r, 3)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	if funcName == "" {
		http.Error(w, "func name not found", http.StatusBadRequest)
		return
	}

	fn, ok := b.fns[funcName]
	if !ok {
		http.Error(w, "funcName '"+funcName+"' not bound", http.StatusBadRequest)
		return
	}

	body, e := mux.ReadBodyString(r)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	in, e := fn.ParseIn(body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	out := reflect.ValueOf(fn.i).Call(in)
	if len(out) == 0 {
		return
	}

	e = checkError(out[len(out)-1])
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	rp, e := json.Marshal(out[0].Interface())
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	mux.WriteJSON(w, rp)
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
