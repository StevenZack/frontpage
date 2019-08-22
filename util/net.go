package util

import "math/rand"

func RandPort() int {
	return rand.Intn(50000) + 1000
}
