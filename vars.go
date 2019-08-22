package frontpage

import (
	"strconv"

	"github.com/StevenZack/frontpage/util"
)

type Vars struct {
	Addr string
}

func NewVars() *Vars {
	return &Vars{
		Addr: "127.0.0.1:" + strconv.Itoa(util.RandPort()),
	}
}
