package frontpage

import "github.com/StevenZack/fasthttp"

type FrontPage struct {
	Router                       *fasthttp.Router
	Port                         string
	chanID                       string
	DisableExitWhenAllPageClosed bool
	wsCounter                    int
	isRunning                    bool
	fnMap                        map[string]Fn
	fnOpen                       func(string)
}
