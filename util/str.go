package util

import (
	"encoding/json"
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
