package frontpage

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"

	"github.com/StevenZack/frontpage/util"
)

type Func struct {
	Name    string
	i       interface{}
	inTypes []reflect.Type
}

func NewFunc(i interface{}) (*Func, error) {
	fn := &Func{
		i: i,
	}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind().String() != "func" {
		return nil, errors.New("NewFunc(): i is not a function")
	}

	fn.Name = runtime.FuncForPC(v.Pointer()).Name()

	for index := 0; index < t.NumIn(); index++ {
		itype := t.In(index)
		fn.inTypes = append(fn.inTypes, itype)
	}
	return fn, nil
}

func (f *Func) ParseIn(body string) ([]reflect.Value, error) {
	strs, e := util.SplitJsonArray(body)
	if e != nil {
		return nil, e
	}

	if len(strs) != len(f.inTypes) {
		return nil, errors.New("invalid input length:" + string(body))
	}

	vs := []reflect.Value{}
	for index, intype := range f.inTypes {
		ptr := reflect.New(intype)
		e := json.Unmarshal([]byte(strs[index]), ptr)
		if e != nil {
			return nil, e
		}

		vs = append(vs, ptr.Elem())
	}

	return vs, nil
}
