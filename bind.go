package frontpage

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func (fp *FrontPage) handleMsg(b []byte) {
	var args []string
	json.Unmarshal(b, &args)
	fn := fp.fnMap[args[0]]
	v := reflect.ValueOf(fn.I)

	in := []reflect.Value{}

	for j := 1; j < len(args); j++ {
		inValue := transformValue(args[j], fn.InsType[j-1])
		in = append(in, inValue)
	}
	fmt.Println("v.call ...:", in)
	v.Call(in)
}

func (fp *FrontPage) BindFunc(f interface{}) {
	fn := transformFn(f)
	fp.fnMap[fn.FnName] = fn
}
