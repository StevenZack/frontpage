package frontpage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/StevenZack/pubsub"

	"github.com/StevenZack/tools/refToolkit"
)

func (fp *FrontPage) handleMsg(b []byte) {
	var args []interface{}
	json.Unmarshal(b, &args)
	fn := fp.fnMap[args[0].(string)]
	v := reflect.ValueOf(fn.I)

	in := []reflect.Value{}

	for j := 1; j < len(args); j++ {
		inValue := transformValue(fmt.Sprint(args[j]), fn.InsType[j-1])
		in = append(in, inValue)
	}
	fmt.Println("v.call ...:", in)
	v.Call(in)
}

func (fp *FrontPage) BindFunc(f interface{}) {
	fn := transformFn(f)
	fp.fnMap[fn.FnName] = fn
}

func transformValue(arg, inType string) reflect.Value {
	switch inType {
	case "int":
		i, _ := strconv.Atoi(arg)
		return reflect.ValueOf(i)
	default:
		return reflect.ValueOf(arg)
	}
}

func transformFn(i interface{}) Fn {
	fn := Fn{}
	t := reflect.TypeOf(i)
	name, e := refToolkit.GetFuncName(i)
	if e != nil {
		panic(e)
	}
	fn.FnName = name
	fn.I = i
	for index := 0; index < t.NumIn(); index++ {
		in := "arg_" + strconv.Itoa(index)
		inType := t.In(index).Name()
		fn.Ins = append(fn.Ins, in)
		fn.InsType = append(fn.InsType, inType)
	}
	return fn
}

func (f *FrontPage) Eval(s string) {
	slice := []string{"eval", s}
	b, e := json.Marshal(slice)
	if e != nil {
		fmt.Println(`eval.marshal error :`, e)
		return
	}
	pubsub.Pub(chanID, string(b))
}

func (f *FrontPage) Invoke(fn string, args ...interface{}) {
	slice := []interface{}{"invoke", fn, args}
	b, e := json.Marshal(slice)
	if e != nil {
		fmt.Println(`Invoke.marshal error :`, e)
		return
	}
	pubsub.Pub(chanID, string(b))
}
