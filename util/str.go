package util

import (
	"encoding/json"
	"reflect"
)

func SplitJsonArray(str string) ([]string, error) {
	vs := []interface{}{}
	e := json.Unmarshal([]byte(str), &vs)
	if e != nil {
		return nil, e
	}
	out := []string{}
	for _, v := range vs {
		b, e := json.Marshal(v)
		if e != nil {
			return nil, e
		}
		out = append(out, string(b))
	}
	return out, nil
}

func RefVsToJson(vs []reflect.Value) (string, error) {
	out := []interface{}{}
	for _, v := range vs {
		out = append(out, v.Interface())
	}

	b, e := json.Marshal(out)
	if e != nil {
		return "", e
	}
	return string(b), nil
}
