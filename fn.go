package frontpage

import (
	"reflect"
	"strconv"

	"github.com/StevenZack/tools/refToolkit"
)

// Fn for Go Template Language
type Fn struct {
	FnName  string
	Ins     []string
	InsType []string
	I       interface{}
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

func transformValue(arg, inType string) reflect.Value {
	switch inType {
	case "int":
		i, _ := strconv.Atoi(arg)
		return reflect.ValueOf(i)
	default:
		return reflect.ValueOf(arg)
	}
}
