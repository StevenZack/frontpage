package util

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(i interface{}) (string, error) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind().String() != "func" {
		return "", errors.New("is not a func")
	}
	name := runtime.FuncForPC(v.Pointer()).Name()
	names := strings.Split(name, ".")
	return names[len(names)-1], nil
}
