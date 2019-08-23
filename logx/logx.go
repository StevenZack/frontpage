package logx

import (
	"github.com/StevenZack/tools/ioToolkit"
)

func Error(args ...interface{}) {
	ioToolkit.Log(args...)
}
